from ast import literal_eval
from dataclasses import dataclass
from pydantic import BaseModel
from typing import Optional, Union, Mapping


@dataclass
class MQMessage:
    data: Union[bytes, dict]
    RqUID: Optional[str] = None
    headers: Optional[dict] = None
    result: Optional[dict] = None

    correlation_id: Optional[str] = None
    reply_to: Optional[str] = None

    async def body(self):
        self.data = literal_eval(self.data.decode('utf-8'))
        return self


@dataclass
class MQHeaders:
    request_headers: Optional[Mapping]

    def __getitem__(self, item):
        return self.request_headers.get(item)

    def get(self, item):
        return self.request_headers.get(item)

    class Config:
        arbitrary_types_allowed = True
        allow_population_by_field_name = True


class HTTPBody(BaseModel):
    data: bytes
    _headers: dict = {'Content-Type': 'application/json'}

    async def body(self):
        return self.data

    @property
    def headers(self, *args, **kwargs):
        new_headers = self._headers
        new_headers.update(**kwargs)
        return new_headers


class UserAuth(BaseModel):
    username: str
    email: Optional[str] = None
    full_name: Optional[str] = None
    disabled: Optional[bool] = None


class UserInDB(UserAuth):
    hashed_password: str
