import traceback

from http import HTTPStatus
from fastapi.exceptions import HTTPException
from grpc import StatusCode
from grpc_interceptor.exceptions import GrpcException


# async def verify_body_request(x_token: str = Header(...)):
#     if x_token != "fake-super-secret-token":
#         raise HTTPException(status_code=400, detail="X-Token header invalid")
#
#
# async def verify_key(x_key: str = Header(...)):
#     if x_key != "fake-super-secret-key":
#         raise HTTPException(status_code=400, detail="X-Key header invalid")
#     return x_key


class ServiceError(HTTPException):

    status_code = HTTPStatus.INTERNAL_SERVER_ERROR

    def __init__(
            self,
            status_code: int = None,
            details: str = None
    ):
        if status_code is not None:
            if status_code == HTTPStatus.OK:
                raise ValueError("The status code for an exception cannot be OK")

            self.status_code = status_code

        HTTPException.__init__(self, status_code=status_code, detail=details, headers={'content-type': 'text/plain'})


class SocketAddressNotFoundError(RuntimeError):
    """ Socket address of a component cannot be found. """
    pass


class ServiceConfigError(RuntimeError):
    """ Component configuration error. """
    pass


class ModelError(ServiceError, GrpcException):
    """
    Any error occurred during model invocation process.
    This error handled by error handlers configured for HTTP and gRPC servers.

    The error message, HTTP and gRPC status codes can be specified.
    """

    def __init__(
            self,
            message: str = None,
            headers: dict = None,
            *,
            http_status: int = HTTPStatus.INTERNAL_SERVER_ERROR,
            grpc_status: StatusCode = StatusCode.INTERNAL
    ):
        """
        Init model error.

        :param message: message displayed in log messages and sent in service response.
        :param http_status: HTTP response status code.
        :param grpc_status: gRPC response status code.
        """

        ServiceError.__init__(self, status_code=http_status, details=message)
        GrpcException.__init__(self, status_code=grpc_status, details=message)

        traceback.print_exc()

        raise HTTPException(
            status_code=http_status,
            detail=message,
            headers=headers
        )
