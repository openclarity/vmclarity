import datetime
import time

from plugin.models import Config, Status
from plugin.scanner.scanner import AbstractScanner

class ExampleScanner(AbstractScanner):
    def __init__(self):
        self.status = Status(state="Ready", message="Started")

    def get_status(self) -> Status:
        return self.status

    def set_status(self, status: Status):
        self.status = status

    async def stop(self, timeout_seconds: int):
        pass

    async def start(self, config: Config):
        self.set_status(Status(state="Running", message="Scan running"))
        time.sleep(10)
        self.set_status(Status(state="Done", message="Scan done"))

    def healthz(self) -> bool:
        return True
