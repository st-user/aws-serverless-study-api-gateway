#########################
# Cognito Identity Pool
#########################

# make cognito identity pool
resource "aws_cognito_identity_pool" "cognito_identity_pool" {
  identity_pool_name               = "my-cognito-identity-pool-name"
  allow_unauthenticated_identities = false

  developer_provider_name = var.cognito_provider_name
}
# Output the identity pool ID
output "cognito_identity_pool_id" {
  value = aws_cognito_identity_pool.cognito_identity_pool.id
}

data "aws_iam_policy_document" "api_gateway_policy" {
  statement {
    effect = "Allow"
    actions = [
      "execute-api:Invoke",
    ]
    resources = ["arn:aws:execute-api:${var.aws_region}:${var.aws_account_id}:*"]
  }
}


resource "aws_iam_policy" "api_gateway_policy" {
  name        = "api_gateway_policy"
  description = "api_gateway_policy"
  policy      = data.aws_iam_policy_document.api_gateway_policy.json
}

data "aws_iam_policy_document" "cognito_policy_document" {
  statement {
    effect = "Allow"
    actions = [
      "mobileanalytics:PutEvents",
      "cognito-identity:*"
    ]
    resources = ["*"]
  }
}

# make IAM policy for "cognito_policy_document"
resource "aws_iam_policy" "cognito_policy" {
  name        = "cognito_policy"
  description = "cognito_policy"
  policy      = data.aws_iam_policy_document.cognito_policy_document.json
}

data "aws_iam_policy_document" "cognito_trust_relationship_policy_document" {
  statement {
    effect = "Allow"
    principals {
      type        = "Federated"
      identifiers = ["cognito-identity.amazonaws.com"]
    }
    actions = [
      "sts:AssumeRoleWithWebIdentity"
    ]
    condition {
      test     = "StringEquals"
      variable = "cognito-identity.amazonaws.com:aud"
      values = [
        "${aws_cognito_identity_pool.cognito_identity_pool.id}"
      ]
    }
    condition {
      test     = "ForAnyValue:StringLike"
      variable = "cognito-identity.amazonaws.com:amr"
      values = [
        "authenticated"
      ]
    }
  }


}

# make IAM role for Cognito
resource "aws_iam_role" "cognito_role" {
  name               = "cognito-authenticated-user-role"
  assume_role_policy = data.aws_iam_policy_document.cognito_trust_relationship_policy_document.json
}

resource "aws_iam_role_policy_attachment" "api_gateway_policy_attachment" {
  role       = aws_iam_role.cognito_role.name
  policy_arn = aws_iam_policy.api_gateway_policy.arn
}

# attach cognito_policy to the IAM role
resource "aws_iam_role_policy_attachment" "cognito_policy_attachment" {
  role       = aws_iam_role.cognito_role.name
  policy_arn = aws_iam_policy.cognito_policy.arn
}

# attach role to the identity pool
resource "aws_cognito_identity_pool_roles_attachment" "cognito_identity_pool_roles_attachment" {
  identity_pool_id = aws_cognito_identity_pool.cognito_identity_pool.id
  roles = {
    "authenticated" = aws_iam_role.cognito_role.arn
  }
}

