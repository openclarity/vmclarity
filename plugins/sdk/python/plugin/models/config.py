from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from plugin.models.base_model import Model
from plugin import util


class Config(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, file=None, input_dir=None, output_file=None, output_schema=None, timeout_seconds=None):  # noqa: E501
        """Config - a model defined in OpenAPI

        :param file: The file of this Config.  # noqa: E501
        :type file: str
        :param input_dir: The input_dir of this Config.  # noqa: E501
        :type input_dir: str
        :param output_file: The output_file of this Config.  # noqa: E501
        :type output_file: str
        :param output_schema: The output_schema of this Config.  # noqa: E501
        :type output_schema: str
        :param timeout_seconds: The timeout_seconds of this Config.  # noqa: E501
        :type timeout_seconds: int
        """
        self.openapi_types = {
            'file': str,
            'input_dir': str,
            'output_file': str,
            'output_schema': str,
            'timeout_seconds': int
        }

        self.attribute_map = {
            'file': 'file',
            'input_dir': 'inputDir',
            'output_file': 'outputFile',
            'output_schema': 'outputSchema',
            'timeout_seconds': 'timeoutSeconds'
        }

        self._file = file
        self._input_dir = input_dir
        self._output_file = output_file
        self._output_schema = output_schema
        self._timeout_seconds = timeout_seconds

    @classmethod
    def from_dict(cls, dikt) -> 'Config':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The Config of this Config.  # noqa: E501
        :rtype: Config
        """
        return util.deserialize_model(dikt, cls)

    @property
    def file(self) -> str:
        """Gets the file of this Config.

        The file with the configuration required by the scanner plugin. This is a path on the filesystem to the config file.   # noqa: E501

        :return: The file of this Config.
        :rtype: str
        """
        return self._file

    @file.setter
    def file(self, file: str):
        """Sets the file of this Config.

        The file with the configuration required by the scanner plugin. This is a path on the filesystem to the config file.   # noqa: E501

        :param file: The file of this Config.
        :type file: str
        """

        self._file = file

    @property
    def input_dir(self) -> str:
        """Gets the input_dir of this Config.

        The directory which should be scanned by the scanner plugin.   # noqa: E501

        :return: The input_dir of this Config.
        :rtype: str
        """
        return self._input_dir

    @input_dir.setter
    def input_dir(self, input_dir: str):
        """Sets the input_dir of this Config.

        The directory which should be scanned by the scanner plugin.   # noqa: E501

        :param input_dir: The input_dir of this Config.
        :type input_dir: str
        """
        if input_dir is None:
            raise ValueError("Invalid value for `input_dir`, must not be `None`")  # noqa: E501

        self._input_dir = input_dir

    @property
    def output_file(self) -> str:
        """Gets the output_file of this Config.

        Path to JSON file where the scanner plugin should store its results.   # noqa: E501

        :return: The output_file of this Config.
        :rtype: str
        """
        return self._output_file

    @output_file.setter
    def output_file(self, output_file: str):
        """Sets the output_file of this Config.

        Path to JSON file where the scanner plugin should store its results.   # noqa: E501

        :param output_file: The output_file of this Config.
        :type output_file: str
        """
        if output_file is None:
            raise ValueError("Invalid value for `output_file`, must not be `None`")  # noqa: E501

        self._output_file = output_file

    @property
    def output_schema(self) -> str:
        """Gets the output_schema of this Config.

        Specifies custom schema the scanner plugin should use to process scan results and save them into `Result.rawData`. Custom schema allows the scanner plugin to be used with third-party tools and services. For example, `cyclondx-json` custom schema can be used to save/parse (JSON) byte stream into/from `Result.rawData` about SBOM findings.  If the custom schema is not supported by the scanner, the scan should fail. When no custom schema is specified, `Result.schema` and `Result.rawData` properties should be empty.   # noqa: E501

        :return: The output_schema of this Config.
        :rtype: str
        """
        return self._output_schema

    @output_schema.setter
    def output_schema(self, output_schema: str):
        """Sets the output_schema of this Config.

        Specifies custom schema the scanner plugin should use to process scan results and save them into `Result.rawData`. Custom schema allows the scanner plugin to be used with third-party tools and services. For example, `cyclondx-json` custom schema can be used to save/parse (JSON) byte stream into/from `Result.rawData` about SBOM findings.  If the custom schema is not supported by the scanner, the scan should fail. When no custom schema is specified, `Result.schema` and `Result.rawData` properties should be empty.   # noqa: E501

        :param output_schema: The output_schema of this Config.
        :type output_schema: str
        """

        self._output_schema = output_schema

    @property
    def timeout_seconds(self) -> int:
        """Gets the timeout_seconds of this Config.

        The maximum time in seconds that a scan started from this scan should run for before being automatically aborted.   # noqa: E501

        :return: The timeout_seconds of this Config.
        :rtype: int
        """
        return self._timeout_seconds

    @timeout_seconds.setter
    def timeout_seconds(self, timeout_seconds: int):
        """Sets the timeout_seconds of this Config.

        The maximum time in seconds that a scan started from this scan should run for before being automatically aborted.   # noqa: E501

        :param timeout_seconds: The timeout_seconds of this Config.
        :type timeout_seconds: int
        """
        if timeout_seconds is None:
            raise ValueError("Invalid value for `timeout_seconds`, must not be `None`")  # noqa: E501
        if timeout_seconds is not None and timeout_seconds < 0:  # noqa: E501
            raise ValueError("Invalid value for `timeout_seconds`, must be a value greater than or equal to `0`")  # noqa: E501

        self._timeout_seconds = timeout_seconds
