#!/usr/bin/env python3

import asyncio

from plugin.models import Config, Status, Metadata, Stop, Result
from plugin.scanner import AbstractScanner
from plugin.server import run_scanner_server, logger


class ExampleScanner(AbstractScanner):
    def __init__(self):
        self.status = Status(state="Ready", message="Scanner ready")

    def get_status(self) -> Status:
        return self.status

    def set_status(self, status: Status):
        self.status = status

    def get_metadata(self) -> Metadata:
        return Metadata(
            name="Example scanner",
            version="v0.1.2",
        )

    async def stop(self, stop: Stop):
        # cleanup logic
        return

    async def start(self, config: Config):
        # Mark scan started
        logger.info("Scanner is running")
        self.set_status(Status(state="Running", message="Scan running"))

        # Example scanning
        await asyncio.sleep(5)
        try:
            result = Result()
            self.export_result(result=result, output_file=config.output_file)
        except Exception as e:
            logger.error(f"Scanner failed with error {e}")
            self.set_status(Status(state="Failed", message="Scan failed"))
            return

        # Mark scan done
        logger.info("Scanner finished running")
        self.set_status(Status(state="Done", message="Scan done"))


if __name__ == '__main__':
    run_scanner_server(ExampleScanner())