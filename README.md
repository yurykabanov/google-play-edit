# Google Play Edit command line tool

This tool provides ability to perform bulk updates of [Edits](https://developers.google.com/android-publisher/edits/)
of Android applications in Google Play market (or at least some part of them).

NOTE: It is in WIP state and can be used only to update listings and upload phone screenshots.

## Quickstart guide

```bash
# 1. Prepare the following:
#    - listings defintions: Language, Title, Short Description, Full Description
#    - (optional) phone screenshots (each language in separate directory)
#
# Example:
#   data/
#     new_listings.json
#     en-US/
#       1.png
#       2.png
#     ru-RU/
#       1.png
#       2.png

# 2. Create and bulk insert updates:
google-play-edit                                        \
    --account ./data/google_service_account.json        \
    --package-name="com.example.my-awesome-application" \
    insert ./data/new_listings.json                     \
    --phone-screenshots ./data/

# 3. Ensure everything is correct (use $id printed by previous step)
google-play-edit                                        \
    --account ./data/google_service_account.json        \
    --package-name="com.example.my-awesome-application" \
    list $id
    
# 4. Commit changes
google-play-edit                                        \
    --account ./data/google_service_account.json        \
    --package-name="com.example.my-awesome-application" \
    commit
```

## Build from scratch

```bash
go mod download
go build -o ./google-play-edit ./cmd/google-play-edit/main.go 
```

## TODO

- More commands (update/delete/verify/etc)
- Upload APKs and other stuff (see [API](https://developers.google.com/android-publisher/api-ref/) for more details)
- Better images support (other types of images)
- Something has to be done with `internal/command` package, the absence of DI makes me sad
- Add more tests
