import os

import requests

from canyantester import canyantester
from click.testing import CliRunner


TARGET = os.environ.get("TARGET", "kamailio:5060")
API_URL = os.environ.get("API_URL", "http://rating-api:8000")


def test_kamailio_call():
    base_dir = os.path.dirname(__file__)
    scenario_file = os.path.join(base_dir, "scenarios", "test_kamailio_call.yaml")
    result = CliRunner().invoke(
        canyantester, ["-a", API_URL, "-t", TARGET, scenario_file]
    )
    assert result.exit_code == 0

    # verify transaction for account 1000
    with requests.post(
        API_URL + "/graphql",
        json={
            "query": """query {
  data:allTransactions(page:1, perPage: 1, sortField:"timestamp_begin", sortOrder:"desc", filter:{
    account_tag:"1000"
  }) {
    tenant
    account_tag
    source
    destination
    inbound
    failed
    failed_reason
    duration
    fee
  }
}"""
        },
    ) as resp:
        data = resp.json()
        assert data == {
            "data": {
                "data": [
                    {
                        "tenant": "default",
                        "account_tag": "1000",
                        "source": "sip:1000@sip.canyan.io",
                        "destination": "sip:1001@sip.canyan.io",
                        "inbound": False,
                        "failed": False,
                        "failed_reason": "",
                        "duration": 1,
                        "fee": 0,
                    }
                ]
            },
            "errors": None,
        }

    # verify transaction for account 1001
    with requests.post(
        API_URL + "/graphql",
        json={
            "query": """query {
  data:allTransactions(page:1, perPage: 1, sortField:"timestamp_begin", sortOrder:"desc", filter:{
    account_tag:"1001"
  }) {
    tenant
    account_tag
    source
    destination
    inbound
    failed
    failed_reason
    duration
    fee
  }
}"""
        },
    ) as resp:
        data = resp.json()
        assert data == {
            "data": {
                "data": [
                    {
                        "tenant": "default",
                        "account_tag": "1001",
                        "source": "sip:1000@sip.canyan.io",
                        "destination": "sip:1001@sip.canyan.io",
                        "inbound": True,
                        "failed": False,
                        "failed_reason": "",
                        "duration": 1,
                        "fee": 0,
                    }
                ]
            },
            "errors": None,
        }
