import json
from abc import abstractmethod, ABCMeta

from framework.schemas.body import Message


class CronModel(metaclass=ABCMeta):
    """
    Example Additional Service Model with endpoint "/cron".
    Service configurable in OpenShift.
    The method must fire events from "/predict"
    """

    @abstractmethod
    async def predict(self, *args, **kwargs):
        return NotImplementedError("Subclasses should implement this")

    @abstractmethod
    async def cron(self, **context):
        return NotImplementedError("Subclasses should implement this")

    async def to_predict(self, item: dict, request_id: str = '42', background_tasks=None):
        """
        Вызывается из пользовательского класса в app.py. Отправляет сообщение в очередь Kafka.
        :param item: пользовательское сообщение
        :param request_id: id запроса
        :param background_tasks: очередель для фоновых задач [fastapi.BackgroundTasks]
        """

        msg = Message(
            data=json.dumps(item).encode('utf-8'),
            headers={'RqUID': request_id},
        )
        result = await self.predict(item=msg, background_tasks=background_tasks)
        return result
