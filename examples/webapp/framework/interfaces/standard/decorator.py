from framework.server.contexts import thread_context
from framework.server.contexts.app_context import shutdown_event
from framework.server.handlers.metric_handler import start_metric, finish_metric
from framework.schemas.exceptions import ModelError
from functools import wraps
from http import HTTPStatus
from typing_extensions import Coroutine


def request_handler(fn):
    """ `request_handler` method decorator.
    For wrap all business requests (/predict, /cron and et.).
    With support request_id tracing and metrics
    """

    @wraps(fn)
    async def wrapper(
            self,
            item,
            RqUID: str = None,
            content_type: str = None,
            *args, **kwargs
    ):
        context_configured = thread_context.contains('request_id')

        if not context_configured:
            token = thread_context.put({
                'request_id': RqUID
            })

        start_time = start_metric()

        _get_data = getattr(self, 'get_data')
        _give_data = getattr(self, 'give_data')

        try:

            receive_data = _get_data(await item.body(), *args)

            if isinstance(receive_data, Coroutine):
                receive_data = await receive_data
            else:
                pass

            result = await fn(
                self, item=receive_data, *args, **kwargs)

        except (SyntaxError, UnicodeDecodeError,) as e:
            finish_metric(start_time, 'failure', 'error')
            raise ModelError(
                message=f'Failed to parse initial data: {e}',
                http_status=HTTPStatus.BAD_REQUEST
            ) from e

        except KeyError as e:
            finish_metric(start_time, 'failure', 'error')
            raise ModelError(
                message=f'Key not found: {e}',
                http_status=HTTPStatus.NOT_FOUND
            )

        except Exception as e:
            finish_metric(start_time, 'failure', 'error')
            raise ModelError(
                message=f'Model is broken: {e}',
                http_status=HTTPStatus.INTERNAL_SERVER_ERROR
            ) from e

        # ||| to interface Model method |||
        try:
            if shutdown_event.is_set():
                status = 'failure'
            else:
                status = 'success'

            finish_metric(start_time, status, 'none')
        except:  # noqa
            pass

        if not context_configured:
            thread_context.remove(token)

        res: dict = await _give_data(
            result, *args, **kwargs
        )

        return res

    return wrapper


def technical_request_handler(fn):
    """ `technical_request_handler` method decorator.
    For wrap technical requests (/health, /info and et.).
    With support request_id tracing but without metrics
    """

    @wraps(fn)
    def wrapper(*args, **kwargs):
        context_configured = thread_context.contains('request_id')

        if not context_configured:
            request_id = kwargs.get('request_id')
            token = thread_context.put({
                'request_id': request_id
            })

        kwargs['shutdown_event'] = shutdown_event

        result = fn(*args, **kwargs)

        if not context_configured:
            thread_context.remove(token)

        return result

    return wrapper
