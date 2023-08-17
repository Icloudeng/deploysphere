# Fixes

> How do I fix "$'\r': command not found" errors running Bash scripts in WSL

Inside WSL:

```bash
sudo apt-get install dos2unix
```

Then,

```bash
dos2unix [file]
```

Full documentation:

```bash
man dos2unix
```

Saved my day, hope it helps.
