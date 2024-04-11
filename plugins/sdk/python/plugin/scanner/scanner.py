from abc import ABC, abstractmethod

from plugin.models.config import Config  # noqa: E501
from plugin.models.status import Status  # noqa: E501
from plugin.models.metadata import Metadata  # noqa: E501

class AbstractScanner(ABC):

    @abstractmethod
    def healthz(self) -> bool:
        pass

    @abstractmethod
    def get_metadata(self) -> Metadata:
        pass

    @abstractmethod
    async def start(self, config: Config):
        pass

    @abstractmethod
    async def stop(self, timeout_seconds: int):
        pass

    @abstractmethod
    def set_status(self, status: Status):
        pass

    @abstractmethod
    def get_status(self) -> Status:
        pass
