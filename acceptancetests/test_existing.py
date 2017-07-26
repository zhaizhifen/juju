#!/usr/bin/env python
"""TODO: add rough description of what is assessed in this module."""

from __future__ import print_function

import argparse
import logging
import sys

from deploy_stack import (
    BootstrapManager,
    )
from utility import (
    add_basic_testing_arguments,
    configure_logging,
    )


__metaclass__ = type


log = logging.getLogger("assess_TEMPLATE")


def assess_TEMPLATE(client):
    # Deploy charms, there are several under ./repository
    client.deploy('cs:ubuntu', num=2)
    # Wait for the deployment to finish.
    client.wait_for_started()
    log.info("TODO: Add log line about any test")
    # TODO: Add specific functional testing actions here.


def parse_args(argv):
    """Parse all arguments."""
    parser = argparse.ArgumentParser(description="TODO: script info")
    # TODO: Add additional positional arguments.
    add_basic_testing_arguments(parser)
    # TODO: Add additional optional arguments.
    return parser.parse_args(argv)


def main(argv=None):
    args = parse_args(argv)
    configure_logging(args.verbose)
    if args.existing:
        bs_manager = BootstrapManager.from_existing(args)
        with bs_manager.existing_context(args.upload_tools,
                                         args.temp_env_name):
                assess_TEMPLATE(bs_manager.client)
    else:
        bs_manager = BootstrapManager.from_args(args)
        with bs_manager.booted_context(args.upload_tools):
                assess_TEMPLATE(bs_manager.client)
    return 0


if __name__ == '__main__':
    sys.exit(main())
