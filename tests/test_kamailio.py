import os

from click.testing import CliRunner
from canyantester import canyantester


TARGET = os.environ.get('TARGET', 'kamailio:5060')


def test_kamailio_call():
    base_dir = os.path.dirname(__file__)
    scenario_file = os.path.join(base_dir, 'scenarios', 'test_kamailio_call.yaml')
    result = CliRunner().invoke(canyantester, ['-t', TARGET, scenario_file])
    assert result.exit_code == 0
