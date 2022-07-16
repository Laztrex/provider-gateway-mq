from framework.server.contexts.singleton import Singleton


def health_checker() -> dict:
    return {"health_status": "running"}


class HealthIndicator(metaclass=Singleton):
    """future"""
    def get_status(self, *args, **kwargs):
        pass
