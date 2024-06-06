# Bud

I am actively working on the project and so far the process to get it running require a few dependencies and to follow some steps in order to run it.

You will need:
 - Postgres Database with Vector modules installed
 - Sqlite
 - Whisper.h
 - Ggml.h

And in case you would like to try Worker rag it also uses:
  - Poppler Utils
  - Docx2txt

In order to run the project you should have the following directories at the src of your project:
  ├── include
  │   ├── ggml.h
  │   └── whisper.h
  ├── lib
  │   └── libwhisper.a
  ├── models
  │   └── ggml-base.en.bin

In order to get them, you must clone the whisper.cpp project and build following the instructions for golang bindings.

  https://github.com/ggerganov/whisper.cpp

You will also need a .env file in the root of the project with the variables:

DBUSER=
DBPASSWORD=
DBHOST=
DBPORT=
