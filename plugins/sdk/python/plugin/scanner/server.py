from flask import Flask, request, jsonify

from plugin.scanner.scanner import AbstractScanner  # noqa: E501
from plugin.models.config import Config  # noqa: E501
from plugin.models.metadata import Metadata  # noqa: E501
from plugin.models.stop import Stop  # noqa: E501

class Server:
    def __init__(self, scanner: AbstractScanner):
        self.app = Flask(__name__)
        self.scanner = scanner
        self.register_routes()

    def start(self, host: str, port: int):
        self.app.run(host=host, port=port)

    def stop(self):
        return jsonify(''), 200

    def register_routes(self):
        self.app.add_url_rule('/healthz', 'get_healthz', self.get_healthz, methods=['GET'])
        self.app.add_url_rule('/metadata', 'get_metadata', self.get_metadata, methods=['GET'])
        self.app.add_url_rule('/config', 'post_config', self.post_config, methods=['POST'])
        self.app.add_url_rule('/status', 'get_status', self.get_status, methods=['GET'])
        self.app.add_url_rule('/stop', 'post_stop', self.post_stop, methods=['POST'])

    def get_healthz(self):
        self.app.logger.info("Received GetHealthz request")
        if self.scanner.healthz():
            return jsonify(''), 200
        else:
            return jsonify(''), 503

    def get_metadata(self):
        self.app.logger.info("Received GetMetadata request")
        metadata = Metadata("1.0")
        return jsonify(metadata.to_dict()), 200

    def post_config(self):
        self.app.logger.info("Received PostConfig request")
        request_data = request.get_json()
        config = Config().from_dict(request_data)

        if self.scanner.get_status().state != "Ready":
            return jsonify({"message": "scanner is not in ready state"}), 409

        self.scanner.start(config)

        return jsonify(''), 201

    def get_status(self):
        self.app.logger.info("Received GetStatus request")
        status = self.scanner.get_status()
        return jsonify(status.to_dict()), 200

    def post_stop(self):
        self.app.logger.info("Received StopScanner request")
        request_data = request.get_json()
        stop_data = Stop().from_dict(request_data)

        self.scanner.stop(stop_data.timeout_seconds)

        return jsonify(''), 201
