# Contemplating the "skyline" interview problem

> Given a 2d description of a set of buildings, compute their skyline profile.

Building description is fairly straight forward, there are some basic choices:
- rectangle defined by two points
- left/right X values with a height

For "the skyline" however, desired outcome becomes less clear:
- Do you just want to trace its outline? If so a list of `<X, Y>` points that
  describe each corner works.
- But if you wanted to fill it in, you might want to directly compute a
  rectangle strip; such strip itself has many representation questions.
- Furthermore, if you were really drawing this stuff (say targeting the GL
  api), you could directly triangulate the skyline (compute a triangle strip).

# Go Setup

See the [go\_setup branch](../../tree/go_setup#go-setup) for a problem harness
done in go; the README in its branch will further point you to solutions.
