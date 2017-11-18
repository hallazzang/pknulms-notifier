import os as os_
import sys
import subprocess

platform_targets = (
    # (os, arch, ext, compression)
    ('darwin', 'amd64', '', 'tar.gz'),  # macOS
    ('linux', 'amd64', '', 'tar.gz'),  # Ubuntu, ...
    ('windows', '386', '.exe', 'zip'),  # Windows 32 bit
    ('windows', 'amd64', '.exe', 'zip'),  # Windows 64 bit
)

for os, arch, ext, compression in platform_targets:
    output_name = f'pknulms-notifier.{os}.{arch}{ext}'
    compress_name = f'pknulms-notifier.{os}.{arch}.{compression}'

    print('target: %s %s' % (os, arch))

    print('building...', end=''); sys.stdout.flush()
    subprocess.run(f'GOOS={os} GOARCH={arch} go build -o build/{output_name}',
        shell=True, stdout=subprocess.DEVNULL)
    print('finished')

    print('deleting old files...', end=''); sys.stdout.flush()
    if os_.path.exists(f'build/{compress_name}'):
        os_.remove(f'build/{compress_name}')
    print('finished')

    print('archiving and compressing...', end=''); sys.stdout.flush()
    if compression == 'tar.gz':
        subprocess.run(f'tar -C build -zvcf build/{compress_name} {output_name}',
            shell=True, stdout=subprocess.DEVNULL)
    elif compression == 'zip':
        subprocess.run(f'zip -j build/{compress_name} build/{output_name}',
            shell=True, stdout=subprocess.DEVNULL)
    print('finished\n')
