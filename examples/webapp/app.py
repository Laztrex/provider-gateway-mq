#!/usr/bin/env python3

import logging
from typing import Union, List, NoReturn
from uuid import UUID

from framework.server.app_server import JOBS, start_server
from framework.server.app_worker import start_model
from framework.schemas.body import MQMessage, MQHeaders

from framework.interfaces.mq.decorator import mq_handler
from framework.interfaces.standard.decorator import request_handler
from framework.interfaces.mq.mq_model import MQModel
from framework.server.handlers.health_check import health_checker

from model.main import start_predict

class Model(MQModel):
    """
    Python -v 3.8+
    Example Service Model
    """

    __slots__ = ('_DATA',)

    async def health(self) -> dict:
        """
        Проверка работоспособности
        Example local request:
        curl http://127.0.0.1:8080/health
        """

        logging.info("Checked health")
        return health_checker()

    async def status_job(self, uid: UUID) -> dict:
        """
        Проверка статуса задачи
        :param uid: Идентификатор задачи
        """

        return JOBS[uid].__dict__

    @mq_handler
    @request_handler
    async def predict(self, item: MQMessage, headers: MQHeaders = None) -> Union[dict, bytes]:
        """
        Основная функция старта расчета модели.

        Example 1 local request:
        Отправка изображения и сохранение ответного изображения "test.png". Подробнее в /model/main.py
        >> curl -X POST -H "Content-Type: media/image" -H "RqUID: a777b633-abf8-4592-a45d-1cc8d19d0d32" -d '{"data": "'"$( base64 -w0 JenkinsBuildStep.PNG )"'"}' http://127.0.0.1:8080/predict >> test.png

        Example 2 local request (background task):
        Сохранение текста запроса
        >> echo -ne '{"background_tasks": "True", "data": "'"$( base64 -w0 IMG_1468.jpeg )"'"}' > request.txt
        >> RESP=$(curl -X POST -H "Content-Type: media/image" -H "RqUID: a777b633-abf8-4592-a45d-1cc8d19d0d32" -d @request.txt http://127.0.0.1:8080/predict)
        Декодирование ответа RESP:
        >> echo $RESP
        {'uid': 'a800b633-abf8-4592-a45d-1cc8d19d0d32', 'status': 'in_progress', 'result': None}

        Запрос на чтение результата задачи
        >> curl -X POST -H "Content-Type: media/image" -H "RqUID: a777b633-abf8-4592-a45d-1cc8d19d0d32" -d '{"status_job": "129868ae-7fb1-48ec-aa55-676606c30205"}' http://127.0.0.1:8080/predict >> test3.png

        :param item: Тело запроса
        :param background_tasks: Очередь фоновых задач
        :param RqUID: Идентификатор запроса
        :param content_type: Формат данных
        :return: Тело ответа
        """

        logging.info("Processing request started")
        logging.info(f"data: {item}")

        ans = await start_model(start_predict, item.data)

        logging.info("Processing request done")
        logging.info(f"correlation_id: {item.correlation_id}")
        item.headers.update({"test": '42'})
        item.result = {"result": ans}

        return item


if __name__ == "__main__":
    start_server(Model)
