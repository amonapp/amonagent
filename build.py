import sys
import os
import subprocess
import time
from datetime import datetime
import shutil
import tempfile
import logging

supported_archs = ["amd64", "i386", "armhf", "arm64"]
BUILD="packaging/build"
PACKAGING="packaging"
AGENT="{0}/amonagent".format(BUILD)

ROOT = os.path.abspath(os.path.dirname(__file__))

def get_version():
    version = run('git describe --always --tags')

    return version


def compile_binary(arch=None):
    version = get_version()
    logging.info("amonagent version: {0}".format(version))

    command = [
        "CGO_ENABLED=0", 
        "GOARCH={0}".format(arch),
        "go build -o amonagent",
        "-ldflags",
        "\"-X main.Version={0}\"".format(version),
        "./cmd/amonagent.go"
    ]

    compile_string = " ".join(command)
    start_time = datetime.utcnow()
    run(compile_string, shell=True)
    end_time = datetime.utcnow()
    total_seconds = (end_time - start_time).total_seconds()
    logging.info("Time taken: {0}s / {1}".format(total_seconds, arch))


def create_package_fs():
    shutil.rmtree(BUILD, ignore_errors=True)
    build_directory = os.path.join(ROOT, BUILD)
    packaging_directory = os.path.join(ROOT, PACKAGING)

    os.makedirs(build_directory)
    os.makedirs(os.path.join(build_directory, "etc", 'opt', 'amonagent'))
    os.makedirs(os.path.join(build_directory, "etc", 'opt', 'amonagent', 'plugins-enabled'))
    os.makedirs(os.path.join(build_directory, 'opt', 'amonagent'))
    os.makedirs(os.path.join(build_directory, "usr", 'bin'))

    binary = os.path.join(ROOT, 'amonagent')

    shutil.copyfile(binary, os.path.join(build_directory, 'opt', 'amonagent', 'amonagent'))
    shutil.copyfile(binary, os.path.join(build_directory, 'usr', 'bin', 'amonagent'))


    os.makedirs(os.path.join(build_directory, "var", 'log', 'amonagent'))
    # os.chmod(os.path.join(build_directory, "var", 'log', 'amonagent'), 755)


    # # /var/run permissions for RPM distros
    os.makedirs(os.path.join(build_directory, "usr", 'lib', 'tmpfiles.d'))
    shutil.copyfile(
        os.path.join(packaging_directory, 'tmpfilesd_amonagent.conf'),
        os.path.join(build_directory, 'usr', 'lib', 'tmpfiles.d', 'amonagent')
    )


    os.makedirs(os.path.join(build_directory, "opt", 'amonagent', 'scripts'))
    shutil.copyfile(
        os.path.join(packaging_directory, 'init.sh'),
        os.path.join(build_directory, 'opt', 'amonagent', 'scripts', 'amonagent.service')
    )

    shutil.copyfile(
        os.path.join(packaging_directory, 'amonagent.service'),
        os.path.join(build_directory, 'opt', 'amonagent', 'scripts', 'amonagent.service')
    )



def run(command, allow_failure=False, shell=False, printOutput=False):
    """
    Run shell command (convenience wrapper around subprocess).
    If printOutput is True then the output is sent to STDOUT and not returned
    """
    out = None
    logging.debug("{}".format(command))
    try:
        cmd = command
        if not shell:
            cmd = command.split()

        stdout = subprocess.PIPE
        stderr = subprocess.STDOUT
        if printOutput:
            stdout = None

        p = subprocess.Popen(cmd, shell=shell, stdout=stdout, stderr=stderr)
        out, _ = p.communicate()
        if out is not None:
            out = out.decode('utf-8').strip()
        if p.returncode != 0:
            if allow_failure:
                logging.warn(u"Command '{}' failed with error: {}".format(command, out))
                return None
            else:
                logging.error(u"Command '{}' failed with error: {}".format(command, out))
                sys.exit(1)
    except OSError as e:
        if allow_failure:
            logging.warn("Command '{}' failed with error: {}".format(command, e))
            return out
        else:
            logging.error("Command '{}' failed with error: {}".format(command, e))
            sys.exit(1)
    else:
        return out

if __name__ == '__main__':
    LOG_LEVEL = logging.INFO
    if '--debug' in sys.argv[1:]:
        LOG_LEVEL = logging.DEBUG
    log_format = '[%(levelname)s] %(funcName)s: %(message)s'
    logging.basicConfig(level=LOG_LEVEL, format=log_format)
    create_package_fs()