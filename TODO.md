- [x] Colored time info.

- [ ] Implement command history. Move between prev/next commands using arrow keys.
  (I guess, to implement this, we have to use some syscalls like `tcgetattr`, and `tcsetattr` to set the terminal to raw mode in which the terminal driver sends the input to the shell without buffering the input (immediately).)

- [ ] `cd` auto-completion.
  (This also requires the terminal driver to be in raw mode.)

- [ ] Pipes

- [ ] Zsh-style directory changing (You can just give a directory like `/usr/bin`, and zsh will `cd` you into that directory.)

- [ ] Redirection (`>`)