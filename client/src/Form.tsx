import { FC, useState, useEffect } from 'react'
import { CognitoIdentityClient, GetCredentialsForIdentityCommand } from '@aws-sdk/client-cognito-identity'
import { SignatureV4 } from '@aws-sdk/signature-v4'
import { Sha256 } from '@aws-crypto/sha256-browser'
import { HttpRequest } from '@aws-sdk/protocol-http'

const TOKEN_ENDPOINT = 'http://localhost:8080/accessInfo'

interface APIGatewayAccessInfoResponse {
	token: TokenResponse
	url: string
}

interface TokenResponse {
	region: string;
	identityPoolId: string;
	identityId: string;
	loginProvider: string;
	token: string;
}

interface APIGatewayResponse {
	message: string
}

const useAPIMessage = (userId: string): { message: string, error: Error | undefined, isLoading: boolean } => {

	const [error, setError] = useState<Error | undefined>(undefined)
	const [isLoading, setIsLoading] = useState<boolean>(false)
	const [message, setMessage] = useState<string>('')

	useEffect(() => {

		const fetchToken = async () => {

			setIsLoading(true)

			try {
				const apiResponse = await fetch(TOKEN_ENDPOINT + `?userId=${userId}`)
				const data = await apiResponse.json() as APIGatewayAccessInfoResponse

				const cognitoClient = new CognitoIdentityClient({ region: data.token.region })
				const getCredentialsCommand = new GetCredentialsForIdentityCommand({
					IdentityId: data.token.identityId,
					Logins: {
						'cognito-identity.amazonaws.com': data.token.token
					},
				})
				const awsResponse = await cognitoClient.send(getCredentialsCommand)


				console.log("awsResponse.Credentials?.AccessKeyId: ", awsResponse.Credentials?.AccessKeyId)
				console.log("awsResponse.Credentials?.SecretKey", awsResponse.Credentials?.SecretKey)
				console.log("awsResponse.Credentials?.SessionToken", awsResponse.Credentials?.SessionToken)

				const url = new URL(`${data.url}/hello`)
				const body = JSON.stringify({
					name: `${userId}-san`
				})
				const signatureV4 = new SignatureV4({
					service: 'execute-api',
					region: data.token.region,
					credentials: {
						accessKeyId: awsResponse.Credentials?.AccessKeyId ?? '',
						secretAccessKey: awsResponse.Credentials?.SecretKey ?? '',
						sessionToken: awsResponse.Credentials?.SessionToken ?? ''
					},
					sha256: Sha256
				})
				const httpRequest = new HttpRequest({
					body,
					headers: {
						'content-type': 'application/json',
						host: url.hostname,
					},
					hostname: url.hostname,
					method: 'POST',
					path: url.pathname
				})
				console.log("url.hostname: ", url.hostname)
				console.log("url.pathname: ", url.pathname.slice(1))

				const signedRequest = await signatureV4.sign(httpRequest)

				console.log("signedRequest: ", signedRequest)

				const responseMessage = await fetch(url, {
					mode: 'cors',
					method: 'POST',
					headers: signedRequest.headers,
					body
				}).then(response => response.json()) as APIGatewayResponse

				setMessage(responseMessage.message)

			} catch (error) {
				setError(error as Error)
			}
			setIsLoading(false)
		}
		fetchToken()

	}, [userId])

	return { message, error, isLoading }
}

export const Form: FC<{ userId: string }> = ({ userId }) => {

	const { message, error, isLoading } = useAPIMessage(userId)

	return (
		<>
			<>{isLoading ? <p>Loading...</p> :
				<>
					{error ? <p>{error.message}</p> : <pre>{message}</pre>}
				</>
			}</>
		</>
	)
}