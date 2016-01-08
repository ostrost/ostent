var child             = require('child_process'),
    gulp              = require('gulp'),
    argv              = require('yargs').argv;
//  reload            = require('gulp-livereload');

var server, makep = {_cmd: 'make'};

function dup(std, proc) {
  return proc[std].on('data', function(d) { process[std].write(d.toString()); });
  return proc[std].on('data', function(d) {
    var lines = d.toString().split('\n');
    for (var i in lines) {
      util.log(lines[i]); // util is gulp-util
    }
  });
}

gulp.task('which gmake', function(cb) {
  var proc = child.spawn('which', ['gmake']);
  proc.stdout.on('end', cb);
  proc.stdout.on('data', function() { makep._cmd = 'gmake'; });
  dup('stderr', proc);
});

gulp.task('gmake dev', ['which gmake'], function() {
  var proc = child.spawn(makep._cmd, '-C . dev'.split(' '));
  dup('stdout', proc);
  dup('stderr', proc);
});

gulp.task('gmake print', ['which gmake'], function(cb) {
  var proc = child.spawn(makep._cmd, 'print-package print-devpackagefiles'.split(' '));
  proc.stdout.on('end', cb);
  proc.stdout.on('data', function(d) {
    var lines = d.toString().split('\n');
    for (var i in lines) {
      var words = d.toString().split('=');
      makep[words[0]] = words[1].replace(/\n$/, '').split(' ');
    }
  });
  dup('stderr', proc);
});

gulp.task('server build', function(cb) {
  var proc = child.spawn('go', ('get -v -race '+makep.package[0]).split(' '));
  proc.stdout.on('end', cb);
  dup('stdout', proc);
  dup('stderr', proc);
});

gulp.task('server run', ['server build'], function() {
  var run = function() {
    var args = [];
    for (var i in argv) {
      if (i == '_' || i == '$0') {
        continue;
      }
      if (i.length == 1) {
        args.push('-'+i);
      } else {
        args.push('--'+i);
      }
      if (typeof(argv[i]) != 'boolean') {
        args.push(argv[i]);
      }
    }
    // console.log('gulp: run args:', args);
    server = child.spawn(makep.package[0].replace(/.*\//, ''), args);
    dup('stdout', server);
    dup('stderr', server);
    // server.stdout.once('data', function() { reload.reload('/'); });
  };
  if (server != null) {
    server.kill();
    setTimeout(run, 1000);
  } else {
    run();
  }
});

gulp.task('server watch', function() {
  gulp.watch(makep.devpackagefiles.map(function(x) { return __dirname + '/' + x; }),
             ['server run']);
});

gulp.task('watch', ['gmake print'], function() {
  // reload.listen();
  gulp.start(['server watch', 'server run']);
  gulp.watch([
    __dirname+'/share/js/*',
    __dirname+'/share/style/*',
    __dirname+'/share/templatesorigin/*'
  ], ['gmake dev']);
});

// Local variables:
// js-indent-level: 2
// js2-basic-offset: 2
// End:
