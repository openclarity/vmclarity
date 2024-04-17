import asyncio

from plugin.models import Config, Status, Metadata
from plugin.scanner.scanner import AbstractScanner
from plugin.consts import API_VERSION


class ExampleScanner(AbstractScanner):
    def __init__(self):
        self.status = Status(state="Ready", message="Scanner ready")

    def get_status(self) -> Status:
        return self.status

    def set_status(self, status: Status):
        self.status = status

    async def stop(self, timeout_seconds: int):
        print("Stop called")
        await asyncio.sleep(timeout_seconds)
        print("Stop done")

    async def start(self, config: Config):
        self.set_status(Status(state="Running", message="Scan running"))
        await asyncio.sleep(10)
        self.set_status(Status(state="Done", message="Scan done"))

    def healthz(self) -> bool:
        return True
    
    def get_metadata(self) -> Metadata:
        return Metadata(
            name="Example scanner",
            version="v0.1.2",
            api_version=API_VERSION,
        )
