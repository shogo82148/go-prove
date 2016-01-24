use File::Temp qw/ tempfile /;
$ENV{GO_PROVE_TEST_HARRIET} ||= do {
    my ($fh, $filename) = tempfile(UNLINK => 1);
    $filename;
};
