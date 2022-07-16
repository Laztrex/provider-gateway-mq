import importlib.machinery
import logging

import yaml
from jinja2 import Template
from os import getenv
from os.path import dirname, join, exists, pardir, abspath

from framework.server.contexts import app_context as context


def _read_metadata(path_metadata: str) -> dict:
    """
    Загрузка файла метаданных (__metadata__) из директории с моделью ../model/
    Нужно для интерфеса Менеджмента модели (/info)
    Подразумевается наличия аттрибута __all__ в файле __metadata__ с содержанием всех необходимых полей для учета
    :param path_metadata: Путь к файлу с метаданными
    :return: Словарь с метаданными
    """

    metadata = {}

    try:

        if exists(path_metadata):
            loaded_metadata = importlib.machinery.SourceFileLoader('__metadata__', path_metadata).load_module()
            all_attrs = getattr(loaded_metadata, '__all__')

            if all_attrs is not None:
                for attr in all_attrs:
                    metadata[attr] = getattr(loaded_metadata, attr)
        return metadata

    except Exception:
        logging.info(f"Failed load __metadata__ from <model>")


def read_config():
    _framework_dir = abspath(join(__file__, pardir, pardir))
    _model_dir = join(_framework_dir, 'model')

    _config_dir = join(dirname(__file__), '')
    _config_file = join(_config_dir, 'resources/app_config.yml')

    _framework_metadata_file = join(_framework_dir, '__metadata__.py')
    _model_metadata_file = join(_model_dir, '__metadata__.py')

    if not exists(_config_file):
        return {}

    with open(_config_file, 'r') as cfg:
        result = yaml.safe_load(
            Template(cfg.read()).render(
                env=getenv,
                app_context_start_time=context.start_time.isoformat(),
                model_metadata=_read_metadata(_model_metadata_file),
                framework_metadata=_read_metadata(_framework_metadata_file)
            )
        )

    return result


CONFIGS = read_config()
