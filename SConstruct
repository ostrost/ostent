# -*- python -*-
import os, os.path, shlex
import SCons

Default(None)

uname = os.uname()
goos, goarch = uname[0].lower(), {'x86_64': 'amd64'}.get(uname[-1], uname[-1])
bindir = 'bin/%s_%s' % (goos, goarch)

Files = lambda ROOT: [os.path.join(sub, f) for sub, _, fs in os.walk(ROOT) for f in fs if not f.startswith('.#')]
def bindata(target, source, env, for_signature):
    fix = source[0].path         if isinstance(source[0], SCons.Node.FS.Dir) else os.path.dirname(source[0].path)
    src = source[0].path +'/...' if isinstance(source[0], SCons.Node.FS.Dir) else os.path.dirname(source[0].path)
    return ' '.join(shlex.split('''
go-bindata
  -pkg    {pkg}
  -o      {o}
  -tags   {tags}
  -prefix {prefix}
          {source}
'''.format(
    o   = target[0],
    pkg = os.path.basename(
        os.path.dirname( # sorry about that
            os.path.abspath(
                target[0].path))),
    prefix = fix,
    source = src,
    tags   = env['TFLAGS'],
)))

def generator(s):
    def generator(target, source, env, for_signature):
        return s.format(target=target, source=source, env=env)
    return generator

go_env = Environment(BUILDERS={'build': Builder(generator=generator('go build -o $TARGETS {source[0]}'))})
go_env['ENV']['PATH'] += ':'+ os.environ['GOROOT'] +'/bin'
go_env['ENV']['GOPATH'] =     os.environ['GOPATH'] +':'+ os.getcwd()

env = Environment(ENV={'PATH': os.environ['PATH'],
                       'HOME': os.path.expanduser('~')}, BUILDERS={
    'bindata': Builder(generator=bindata),
    'sass':    Builder(action='sass $SOURCES $TARGETS'),
    'jsx':     Builder(action='jsx <$SOURCES  >/dev/null && jsx <$SOURCES 2>/dev/null >$TARGETS'),
    'amberpp': Builder(generator=generator('{source[2]} -defines {source[0]} $FLAG -output $TARGET {source[1]}')),
})

assets    = (Dir('assets/'), Files('assets/'))
templates = ('templates.html/index.html',
             'templates.html/usepercent.html',
             'templates.html/tooltipable.html',
             Files('templates.html/'))
Default(env.Clone(TFLAGS= 'production')       .bindata('src/ostential/view/bindata.production.go',   source=templates))
Default(env.Clone(TFLAGS='!production -debug').bindata('src/ostential/view/bindata.devel.go',        source=templates))
Default(env.Clone(TFLAGS= 'production')       .bindata('src/ostential/assets/bindata.production.go', source=assets))
Default(env.Clone(TFLAGS='!production -debug').bindata('src/ostential/assets/bindata.devel.go',      source=assets))

Default(env.sass('assets/css/index.css', 'style/index.scss'))

amberpp = go_env.build('%s/amberpp' % bindir, source=(Dir('amberp/amberpp'), Glob('src/amberp/amberpp/*.go')))
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

jscript_jsx = env.Clone(FLAG='-j').amberpp(
    'tmp/jscript.jsx',
    ('amber.templates/defines.amber',
     'amber.templates/jscript.amber',
     amberpp))
Default(jscript_jsx)
Default(env.jsx(target='assets/js/gen/jscript.js', source=jscript_jsx))

build_env = Environment(ENV={'PATH': os.environ['PATH'], 'GOPATH': os.environ['GOPATH']}, # +':'+ os.getcwd()
                        BUILDERS={'build': Builder(generator=generator('go build -tags production -o $TARGET {source[0]}'))})
# non-Default
ostent = build_env.build('%s/ostent' % bindir, (
    Dir(    'ostent'    ), # <- package name
    Dir('src/ostent'),
    Dir('src/ostential'),
        'src/ostential/view/bindata.production.go',
        'src/ostential/view/bindata.devel.go',
        'src/ostential/assets/bindata.production.go',
        'src/ostential/assets/bindata.devel.go',
))
Alias('build', ostent)
