from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from plugin.models.base_model import Model
from plugin import util


class Status(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, state=None, message=None, last_transition_time=None):  # noqa: E501
        """Status - a model defined in OpenAPI

        :param state: The state of this Status.  # noqa: E501
        :type state: str
        :param message: The message of this Status.  # noqa: E501
        :type message: str
        :param last_transition_time: The last_transition_time of this Status.  # noqa: E501
        :type last_transition_time: datetime
        """
        self.openapi_types = {
            'state': str,
            'message': str,
            'last_transition_time': datetime
        }

        self.attribute_map = {
            'state': 'state',
            'message': 'message',
            'last_transition_time': 'lastTransitionTime'
        }

        self._state = state
        self._message = message
        self._last_transition_time = last_transition_time

    @classmethod
    def from_dict(cls, dikt) -> 'Status':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The Status of this Status.  # noqa: E501
        :rtype: Status
        """
        return util.deserialize_model(dikt, cls)

    @property
    def state(self) -> str:
        """Gets the state of this Status.

        Describes the status of scanner. | Status         | Description                                                   | | -------------- | ------------------------------------------------------------- | | NotReady       | Initial state when the scanner container starts               | | Ready          | Scanner setup is complete and it is ready to receive requests | | Running        | Scanner config was received and the scanner is running        | | Failed         | Scanner failed                                                | | Done           | Scanner is completed successfully                             |   # noqa: E501

        :return: The state of this Status.
        :rtype: str
        """
        return self._state

    @state.setter
    def state(self, state: str):
        """Sets the state of this Status.

        Describes the status of scanner. | Status         | Description                                                   | | -------------- | ------------------------------------------------------------- | | NotReady       | Initial state when the scanner container starts               | | Ready          | Scanner setup is complete and it is ready to receive requests | | Running        | Scanner config was received and the scanner is running        | | Failed         | Scanner failed                                                | | Done           | Scanner is completed successfully                             |   # noqa: E501

        :param state: The state of this Status.
        :type state: str
        """
        allowed_values = ["NotReady", "Ready", "Running", "Failed", "Done"]  # noqa: E501
        if state not in allowed_values:
            raise ValueError(
                "Invalid value for `state` ({0}), must be one of {1}"
                .format(state, allowed_values)
            )

        self._state = state

    @property
    def message(self) -> str:
        """Gets the message of this Status.

        Human readable message.  # noqa: E501

        :return: The message of this Status.
        :rtype: str
        """
        return self._message

    @message.setter
    def message(self, message: str):
        """Sets the message of this Status.

        Human readable message.  # noqa: E501

        :param message: The message of this Status.
        :type message: str
        """

        self._message = message

    @property
    def last_transition_time(self) -> datetime:
        """Gets the last_transition_time of this Status.

        Last date time when the status has changed.  # noqa: E501

        :return: The last_transition_time of this Status.
        :rtype: datetime
        """
        return self._last_transition_time

    @last_transition_time.setter
    def last_transition_time(self, last_transition_time: datetime):
        """Sets the last_transition_time of this Status.

        Last date time when the status has changed.  # noqa: E501

        :param last_transition_time: The last_transition_time of this Status.
        :type last_transition_time: datetime
        """
        if last_transition_time is None:
            raise ValueError("Invalid value for `last_transition_time`, must not be `None`")  # noqa: E501

        self._last_transition_time = last_transition_time
