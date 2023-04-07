- [ ] Refactoring (Always a possibility)

- [x] Colored time info.

- [ ] Implement command history.
  - [ ] I guess, to implement this, we have to use some syscalls like `tcgetattr`, and `tcsetattr` to set the terminal to raw mode in which the terminal driver sends the input to the shell without buffering the input (immediately).

- [ ] `cd` auto-completion.
  - [ ] This also requires the terminal driver to be in raw mode.

- [ ] A scripting language (.gsh).

- [ ] Pipes

- [ ] Optimize search for executables in `$PATH` using caches (It is an unnecessary complexity. But, maybe...)