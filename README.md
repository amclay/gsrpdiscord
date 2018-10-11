# gsrpdiscord

This bot will automatically update a membership role for a user if they have their game set to ARK (or GSRP/Gunsmoke) and check their Twitch title to ensure that they have GSRP or gunsmoke in it.

To add the bot, have server owner go this URL
https://discordapp.com/oauth2/authorize?client_id=498388551277215745&scope=bot

## Required Setup

1. Add bot (above)
2. Add bot to a "Bot" role in server
3. Allow this role access to change roles
4. Reorder this role above any roles that need changing (it can't alter roles of users that are above it)
5. Create "streamer" role
6. Ensure this is below the bot role so it can alter users in this role
