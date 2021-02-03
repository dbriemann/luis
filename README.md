# luis : LUcid Image Sharing

A simple, easy to use, self-hosted, image sharing solution for family or friends.

# Documentation

## ENV vars

The following environment variables need to be set for luis:

- `LUIS_ADMIN_EMAIL`: the email(login) of the admin account.

# Dependencies

The following dependencies need to be installed on the server running Luis:

- imagemagick (for exif extraction and image manipulation)

# TODO Features

- Users with different rights
- Users can see which images are new
- Create/delete albums
- Add/remove pictures/videos
- After upload create thumb+web version of original picture
- Allow download of single pictures and ZIPs of whole albums
- Create meta data for every picture (read exif)
- Sort by meta data e.g. date taken
