from abc import ABCMeta
from typing import Callable

from fastapi_utils.tasks import repeat_every


class RepeatModel(metaclass=ABCMeta):
    """
    Example Additional Service Model with endpoint "/cron".
    Service configurable in current web-serving.
    The method must fire events from "/predict"
    """

    @staticmethod
    def add_function(function: Callable):
        setattr(RepeatModel, 'repeat_function', function)

    async def repeat(self):
        await self._repeat()

    @staticmethod
    @repeat_every(seconds=5, raise_exceptions=True)
    async def _repeat(*args, **kwargs):
        return getattr(RepeatModel, 'repeat_function')(*args, **kwargs)
