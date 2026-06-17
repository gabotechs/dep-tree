import 'dart:async';
import 'package:sample_app/foo.dart';
import 'package:other_pkg/other.dart';
import 'utils/helper.dart' show helper;
import 'src/bar.dart' as bar;

export 'src/bar.dart' show Bar;

part 'main.g.dart';

void main() {
  print(foo());
  helper();
  bar.Bar();
}
