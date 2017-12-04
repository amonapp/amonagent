import sys
import os
import subprocess
import time
from datetime import datetime
import shutil
import glob
import tempfile
import logging

supported_archs = ["amd64", "i386", "armhf", "arm64"]
BUILD="packaging/build"
PACKAGING="packaging"
AGENT="{0}/amonagent".format(BUILD)

DEV_DEPENDECIES = ['reprepro', 'createrepo']


PACKAGES_PATH="/home/martin/amon/amon-packages"
DEBIAN_REPO_PATH="{0}/debian/".format(PACKAGES_PATH)
RPM_REPO_PATH="{0}/centos/".format(PACKAGES_PATH)

ROOT = os.path.abspath(os.path.dirname(__file__))

def get_version():
    version = run('git describe --always --tags')

    return version


def compile_binary(arch=None):
    os.remove(os.path.join(ROOT, 'amonagent'))
    version = get_version()
    logging.info("amonagent version: {0}".format(version))
    logging.info("Compiling binary for: {0}".format(arch))

    additional_env_params = []
    if arch == "i386":
        arch = "386"
    if arch == "armhf":
        arch = 'arm'
        additional_env_params.append("GOARM=6")
    if arch == "arm64":
        additional_env_params.append("GOARM=7")

    command = [
        "CGO_ENABLED=0", 
        "GOARCH={0}".format(arch),
        "go build -o amonagent",
        "-ldflags",
        '"-X main.Version={0}"'.format(version),
        "./cmd/amonagent.go"
    ]

    command = additional_env_params + command
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


def fpm_build(arch=None, output=None):
    logging.info("Building package for: {0} / {1}".format(arch, output))
    build_directory = os.path.join(ROOT, BUILD)
    packaging_directory = os.path.join(ROOT, PACKAGING)
    
    command = [
        'fpm',
        '--epoch 1',
        '--force',
        '--input-type dir',
        '--output-type {0}'.format(output),
        '--chdir {0}'.format(build_directory),
        '--maintainer "Amon Packages <packages@amon.cx>"',
        '--url "http://amon.cx/"',
        '--description "Amon monitoring agent"',
        '--version {0}'.format(get_version()),
        '--conflicts "amonagent < {0}"'.format(get_version()), 
        '--vendor Amon',
        '--name amonagent',
        '--depends "adduser"',
        '--depends "sysstat"',
        '--architecture "{0}"'.format(arch),
        '--post-install {0}'.format(os.path.join(packaging_directory, 'postinst.sh')),
        '--post-uninstall {0}'.format(os.path.join(packaging_directory, 'postrm.sh')),
        '--pre-uninstall {0}'.format(os.path.join(packaging_directory, 'prerm.sh')),
    ]

    command_string = " ".join(command)
    run(command_string, shell=True)


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

def update_repositories():
    for file in glob.glob("*.deb"):
        logging.info("Copying {0} to {1}".format(file, DEBIAN_REPO_PATH))
        shutil.copy(file, DEBIAN_REPO_PATH)

    command = [
        'find',
        '{0}'.format(DEBIAN_REPO_PATH),
        '-name "*.deb"',
        '-exec reprepro',
        '--ask-passphrase',
        '-Vb',
        '{0}repo'.format(DEBIAN_REPO_PATH),
        'includedeb amon {} \;',
    ]

    command_string = " ".join(command)
    run(command_string, shell=True)

    for file in glob.glob("*.rpm"):
        logging.info("Copying {0} to {1}".format(file, RPM_REPO_PATH))
        shutil.copy(file, RPM_REPO_PATH)


    command = [
        'createrepo',
        '--update {0}'.format(RPM_REPO_PATH),
    ]

    command_string = " ".join(command)
    run(command_string, shell=True)

def upload():
    run("sudo ntpdate -u pool.ntp.org", shell=True)

    command = [
        "aws s3 sync",
        "{0}/debian/repo".format(PACKAGES_PATH),
        "s3://packages.amon.cx/repo/",
        "--region=eu-west-1",
        "--profile=personal",
    ]

    command_string = " ".join(command)
    run(command_string, shell=True)  
    
    command = [
        "aws s3 sync",
        "{0}/centos".format(PACKAGES_PATH),
        "s3://packages.amon.cx/rpm/",
        "--region=eu-west-1",
        "--profile=personal",
    ]
    
    command_string = " ".join(command)
    run(command_string, shell=True)  

def cleanup():
    for file in glob.glob("*.rpm"):
        logging.info("Cleaning up {0}".format(file))
        os.remove(file)
    
    for file in glob.glob("*.deb"):
        logging.info("Cleaning up {0}".format(file))
        os.remove(file)

if __name__ == '__main__':
    LOG_LEVEL = logging.INFO
    if '--debug' in sys.argv[1:]:
        LOG_LEVEL = logging.DEBUG
    log_format = '[%(levelname)s] %(funcName)s: %(message)s'
    logging.basicConfig(level=LOG_LEVEL, format=log_format)


    for arch in supported_archs:
        compile_binary(arch=arch)
        create_package_fs()
        fpm_build(arch=arch, output='rpm')
        fpm_build(arch=arch, output='deb')

    update_repositories()
    upload()
    cleanup()