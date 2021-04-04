load("assert.star", "assert")

def test_main():
    assert.fails(lambda: filter(0, []), "want callable")
    assert.fails(lambda: filter([], 10), "want callable")
    assert.fails(lambda: filter(lambda x: x, 10), "want iterable")
    assert.fails(lambda: filter(lambda x: x, [], 10), "want exactly 2")
    assert.fails(lambda: filter(lambda x: x, [], 10, abc="abc"), "does not accept keyword arguments")

    tests = [
        lambda x: x % 2 == 0,
    ]
    for test in tests:
        assert.eq(
            list(filter(test, [1, 2, 3, 4])),
            [2, 4],
        )
        assert.eq(
            list(filter(test, range(1, 5))),
            [2, 4],
        )

    # Filter empty.
    tests = [
        None,
        bool,
        len,
        lambda item: item,
    ]
    for test in tests:
        assert.eq(
            list(filter(test, ["abc", "", "abcd", ""])),
            ["abc", "abcd"],
        )

    # Keep all.
    tests = [
        lambda x: True,
        lambda x: 1,
    ]
    for test in tests:
        assert.eq(
            list(filter(test, [1, 2, 3])),
            [1, 2, 3],
        )
        assert.eq(
            list(filter(test, range(1, 4))),
            [1, 2, 3],
        )

    # Filter all.
    tests = [
        lambda x: False,
        lambda x: None,
        lambda x: 0,
    ]
    for test in tests:
        assert.eq(
            list(filter(test, [1, 2, 3])),
            [],
        )
        assert.eq(
            list(filter(test, range(1, 4))),
            [],
        )

test_main()