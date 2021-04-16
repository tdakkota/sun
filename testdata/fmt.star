load("assert.star", "assert")

def test_main():
    assert.fails(lambda: hex([]), "want int")
    assert.fails(lambda: oct([]), "want int")
    assert.fails(lambda: bin([]), "want int")

    tests = [(-3, '-0x3'), (-2, '-0x2'), (-1, '-0x1'), (0, '0x0'), (1, '0x1'), (2, '0x2')]
    for test in tests:
        assert.eq(
            hex(test[0]),
            test[1],
        )

    tests = [(-3, '-0o3'), (-2, '-0o2'), (-1, '-0o1'), (0, '0o0'), (1, '0o1'), (2, '0o2')]
    for test in tests:
        assert.eq(
            oct(test[0]),
            test[1],
        )

    tests = [(-3, '-0b11'), (-2, '-0b10'), (-1, '-0b1'), (0, '0b0'), (1, '0b1'), (2, '0b10')]
    for test in tests:
        assert.eq(
            bin(test[0]),
            test[1],
        )

test_main()