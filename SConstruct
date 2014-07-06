# -*- python -*-
import os, os.path, shlex, re
import SCons

Default(None)

uname = os.uname()
goos, goarch = uname[0].lower(), {'x86_64': 'amd64'}.get(uname[-1], uname[-1])
bindir = 'bin/%s_%s' % (goos, goarch)

# Files = lambda ROOT: [os.path.join(sub, f) for sub, _, fs in os.walk(ROOT) for f in fs if not f.startswith('.#')]
class Files(list):
    def __init__(self, ROOT, IGNORX=None):
        ignorx = self.IGNORX = IGNORX
        rx = re.compile(ignorx) if ignorx is not None else None

        files = []
        for sub, _, fs in os.walk(ROOT):
            for f in fs:
                if f.startswith('.#'):
                    continue
                p = os.path.join(sub, f)
                if rx is None or rx.match(p) is None:
                    files += p,
        super(Files, self).__init__(files)

def bindata(target, source, env, for_signature):
    fix = source[0].path         if isinstance(source[0], SCons.Node.FS.Dir) else os.path.dirname(source[0].path)
    src = source[0].path +'/...' if isinstance(source[0], SCons.Node.FS.Dir) else os.path.dirname(source[0].path)
    return ' '.join(shlex.split('''
go-bindata
  -pkg    {pkg}
  -o      {target[0]}
  -tags   $FLAVOR
  -prefix {prefix}
  $IGNORE
          {src}
'''.format(
    pkg = os.path.basename(
        os.path.dirname( # sorry about that
            os.path.abspath(
                target[0].path))),
    prefix = fix,
    src    = src,
    target = target,
)))

def generator(s):
    def generator(target, source, env, for_signature):
        return s.format(target=target, source=source, env=env)
    return generator

go = Environment(
    BUILDERS={
        'build': Builder(generator=generator('go build $TAGSARGS -o $TARGETS {source[0]}'))
    }, ENV={
          'PATH': os.environ[  'PATH'] +':'+ os.environ['GOROOT'] +'/bin',
        'GOPATH': os.environ['GOPATH'] +':'+ os.getcwd()
    })

env = Environment(
    ENV={
        'PATH': os.environ['PATH'] +':'+ os.getcwd() + '/node_modules/.bin',
        'HOME': os.path.expanduser('~')
    }, BUILDERS={
        'bindata': Builder(generator=bindata),
        'sass':    Builder(action='sass $SOURCES $TARGETS'),
        'jsx':     Builder(action='jsx <$SOURCES  >/dev/null && jsx <$SOURCES 2>/dev/null >$TARGETS'),
        'coffee':  Builder(action='coffee -p $SOURCES >/dev/null && coffee -o $TARGETS.dir $SOURCES'),
        'amberpp': Builder(generator=generator('{source[2]} -defines {source[0]} $MODE -output $TARGET {source[1]}')),
    })

assets    = (Dir('assets/'), Files('assets/')) # , IGNORX='assets/js/bundle'))
templates = ('templates.html/index.html',
             'templates.html/usepercent.html',
             'templates.html/tooltipable.html',
             Files('templates.html/'))

Default(env.Clone(FLAVOR= 'production')       .bindata(source=templates, target='src/ostential/view/bindata.production.go'))
Default(env.Clone(FLAVOR='!production -debug').bindata(source=templates, target='src/ostential/view/bindata.devel.go'))

Default(env.Clone(FLAVOR= 'production', IGNORE=('-ignore '+ assets[1].IGNORX if assets[1].IGNORX is not None else None))
        .bindata(source=assets, target='src/ostential/assets/bindata.production.go'))
Default(env.Clone(FLAVOR='!production -debug')
        .bindata(source=assets, target='src/ostential/assets/bindata.devel.go'))

Default(env.sass('assets/css/index.css', 'style/index.scss'))

amberpp = go.build('%s/amberpp' % bindir, source=(Dir('amberp/amberpp'), Glob('src/amberp/amberpp/*.go')))
Default(amberpp)

Default(env.amberpp(
    'templates.html/index.html',
    ('amber.templates/defines.amber',
     'amber.templates/index.amber',
     amberpp)))

Default(env.amberpp(
    'templates.html/usepercent.html',
    ('amber.templates/defines.amber',
     'amber.templates/usepercent.amber',
     amberpp)))

Default(env.amberpp(
    'templates.html/tooltipable.html',
    ('amber.templates/defines.amber',
     'amber.templates/tooltipable.amber',
     amberpp)))

jscript_jsx = env.Clone(MODE='-j').amberpp(
    'tmp/jscript.jsx',
    ('amber.templates/defines.amber',
     'amber.templates/jscript.amber',
     amberpp))
Default(jscript_jsx)
Default(env.jsx(target='assets/js/gen/jscript.js', source=jscript_jsx))

Default(env.coffee(target='assets/js/milk/index.js', source='coffee/index.coffee'))

# non-Default
ostent = go.Clone(TAGSARGS='-tags production').build('%s/ostent' % bindir, (
    Dir(    'ostent'), # <- package name
    Dir('src/ostent'),
    Dir('src/ostential'),
        'src/ostential/view/bindata.production.go',
        'src/ostential/view/bindata.devel.go',
        'src/ostential/assets/bindata.production.go',
        'src/ostential/assets/bindata.devel.go',
))
Alias('b',     ostent)
Alias('build', ostent)
