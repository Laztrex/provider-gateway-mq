from logging import config, Filter, LogRecord
from os import getenv
from os.path import join, dirname

from jinja2 import Template
from yaml import safe_load

_config_dir = join(dirname(__file__), '')
_config_file = join(_config_dir, 'resources/logger_config.yml')


class HealthCheckFilter(Filter):
    def filter(self, record: LogRecord) -> bool:
        return record.getMessage().find("/health") == -1


def load(disable_exist_logs: bool = True):
    with open(_config_file, 'r') as file:
        config.dictConfig(safe_load(Template(file.read()).render(env=getenv, disable_exist=disable_exist_logs)))
