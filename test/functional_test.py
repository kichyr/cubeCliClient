"""
This module contains functional tests of the cubeclient tool.
"""
import subprocess


class TestCubeCli:
    """
    TestCubeCli provide functional test of cubeclient app.
    This tests use test server that up in setup_class method
    and shutdown when all tests done.
    All this tests use binaries for client and server.
    """
    test_srv_port = 8091
    srv_p = None

    def setup_class(self):
        self.srv_p = subprocess.Popen(
            ["./testserver", str(self.test_srv_port)])

    def teardown_class(self):
        self.srv_p.kill()

    def test_bad_args_number(self):
        try:
            subprocess.check_output(
                ['./cubeclient', 'localhost']
            )
        except subprocess.CalledProcessError as exc:
            assert b"wrong number of cli args" in exc.output

    def test_bad_port_format(self):
        try:
            subprocess.check_output(
                ['./cubeclient', 'localhost', "port", 'test1', 'write']
            )
        except subprocess.CalledProcessError as exc:
            print(exc.output)
            assert b"wrong format of port" in exc.output

    def test_bad_server_host(self):
        try:
            subprocess.check_output(
                ['./cubeclient', 'localhost', "999999", 'test1', 'write']
            )
        except subprocess.CalledProcessError as exc:
            assert b"Error during executing request" in exc.output

    def test_ok_response(self):
        output = subprocess.check_output(
            [
                './cubeclient', 'localhost',
                str(self.test_srv_port), 'test1', 'write'
            ]
        )
        assert(
            output ==
            b"client_id: test_client_id \n" +
            b"client_type: 2002 \n" +
            b"expires_in: 3600 \n" +
            b"user_id: 0 \n" +
            b"username: testuser0@mail.ru \n"
        )

    def test_bad_scope_response(self):
        output = subprocess.check_output(
            [
                './cubeclient', 'localhost',
                str(self.test_srv_port), 'test1', 'admin'
            ]
        )
        assert(
            output ==
            b"error: 6 \n" +
            b"message: CUBE_OAUTH2_ERR_BAD_SCOPE \n"
        )
