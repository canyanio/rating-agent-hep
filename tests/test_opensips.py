import os

import requests
import time

from canyantester import canyantester
from click.testing import CliRunner


TARGET = os.environ.get("TARGET_OPENSIPS", "opensips:5060")
API_URL = os.environ.get("API_URL", "http://rating-api:8000")


def test_opensips_call():
    base_dir = os.path.dirname(__file__)
    scenario_file = os.path.join(base_dir, "scenarios", "test_opensips_call.yaml")
    result = CliRunner().invoke(
        canyantester, ["--verbose", "-a", API_URL, "-t", TARGET, scenario_file]
    )
    assert result.exit_code == 0