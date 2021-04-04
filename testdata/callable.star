load("assert.star", "assert")

def function():
    return None

def test_main():
    assert.fails(lambda: callable(lambda x: x, 10), "want exactly 1")
    assert.fails(lambda: callable(lambda x: x, abc="abc"), "does not accept keyword arguments")

    # True
    tests = [
        lambda x: x % 2 == 0,
        function,
    ]
    for test in tests:
        assert.true(
            callable(test),
        )

    # False
    tests = [
        None,
        "",
        1,
        0,
        10,
        range(10),
        True,
        False
    ]
    for test in tests:
        assert.true(
            not callable(test),
        )

test_main()