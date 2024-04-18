from flask import Flask, request, jsonify, copy_current_request_context
from threading import Thread
import logging
import sys
import asyncio

from plugin import encoder
from plugin.scanner.scanner import AbstractScanner  # noqa: E501
from plugin.models.config import Config  # noqa: E501
from plugin.models.error_response import ErrorResponse  # noqa: E501
from plugin.models.stop import Stop  # noqa: E501
from plugin.scanner.config import Config


class Server:
    def __init__(self, scanner: AbstractScanner):
        self.app = Flask(__name__)
        self.app.json_encoder = encoder
        self.scanner = scanner
        self.config = Config()
        self.register_routes()

        logging.basicConfig(stream=sys.stdout,
                            level=logging.getLevelName(self.config.log_level.upper()))

    def start(self):
        self.app.run(host=self.config.get_host(), port=self.config.get_port())

    def stop(self):
        return "", 200

    def register_routes(self):
        self.app.add_url_rule('/healthz', 'get_healthz', self.get_healthz, methods=['GET'])
        self.app.add_url_rule('/metadata', 'get_metadata', self.get_metadata, methods=['GET'])
        self.app.add_url_rule('/config', 'post_config', self.post_config, methods=['POST'])
        self.app.add_url_rule('/status', 'get_status', self.get_status, methods=['GET'])
        self.app.add_url_rule('/stop', 'post_stop', self.post_stop, methods=['POST'])

    def get_healthz(self):
        self.app.logger.info("Received GetHealthz request")
        if self.scanner.healthz():
            return "", 200
        else:
            return "", 503

    def get_metadata(self):
        self.app.logger.info("Received GetMetadata request")

        return self.scanner.get_metadata(), 200

    def post_config(self):
        self.app.logger.info("Received PostConfig request")
        request_data = request.get_json()
        config = Config().from_dict(request_data)

        if self.scanner.get_status().state != "Ready":
            return ErrorResponse(message="scanner is not in ready state"), 409

        @copy_current_request_context
        def start_scanner(config):
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
            try:
                loop.run_until_complete(self.scanner.start(config))
            finally:
                loop.close()

        Thread(target=start_scanner, args=(config,)).start()

        return "", 201

    def get_status(self):
        self.app.logger.info("Received GetStatus request")
        status = self.scanner.get_status()
        return status, 200

    def post_stop(self):
        self.app.logger.info("Received StopScanner request")
        request_data = request.get_json()
        stop_data = Stop().from_dict(request_data)

        @copy_current_request_context
        def stop_scanner(stop_data):
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
            try:
                loop.run_until_complete(self.scanner.stop(stop_data.timeout_seconds))
            finally:
                loop.close()

        Thread(target=stop_scanner, args=(stop_data,)).start()

        return "", 201
