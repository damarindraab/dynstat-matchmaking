# Deployment

After completing testing, deploy the app to AccelByte Gaming Services.

## 1. Create An Extend Override App

If you do not already have one, create a new [Extend Override App](https://docs.accelbyte.io/gaming-services/services/extend/override/matchmaking/get-started-matchmaking-v2/#create-the-extend-app).

On the App Detail page, take note of:

- `Namespace`
- `App Name`

Under Environment Configuration, set these secrets:

- `AB_CLIENT_ID`
- `AB_CLIENT_SECRET`

## 2. Build And Push The Container Image

Use [extend-helper-cli](https://github.com/AccelByte/extend-helper-cli) to build and upload the image:

```
extend-helper-cli image-upload --login --namespace <namespace> --app <app-name> --image-tag v0.0.1
```

> :warning: Run this command from your project directory. If you are in a different directory, add the `--work-dir <project-dir>` option to specify the correct path.

## 3. Deploy The Image

On the App Detail page:

- Click Image Version History
- Select the image you just pushed
- Click Deploy Image

## Customization Tips

1. **Add more stats**: Update the `statistics` array in your rules JSON with any stat codes.
2. **Change the enriched key**: Modify `enriched_key` to match your AGS matching rules attribute.
3. **Implement custom MakeMatches**: If you need custom matching logic beyond AGS defaults, implement `MakeMatches` in `pkg/server/matchmaker.go`.
4. **Add validation rules**: Extend `ValidateTicket` with additional checks (e.g., minimum values, restricted stats).

For more details, see the [AccelByte documentation](https://docs.accelbyte.io/gaming-services/services/extend/override/matchmaking/customization-matchmaking-v2/).
