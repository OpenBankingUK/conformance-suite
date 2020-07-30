# Swagger patch notes

## REFAPP-1083-Transactionsubcode
Changes to remove banktransaction, code and subcode length checking 
- check incorrectly limited field lengths to 4 chars.

```sh

---
 .../account-info-swagger-flattened.json       | 24 +++++-------------
 .../account-info-swagger-flattened.json       | 24 +++++-------------
 .../account-info-swagger-flattened.json       | 25 ++++++-------------
 3 files changed, 19 insertions(+), 54 deletions(-)

diff --git a/pkg/schema/spec/v3.1.3/account-info-swagger-flattened.json b/pkg/schema/spec/v3.1.3/account-info-swagger-flattened.json
index 5a9895dc..ffe557cf 100644
--- a/pkg/schema/spec/v3.1.3/account-info-swagger-flattened.json
+++ b/pkg/schema/spec/v3.1.3/account-info-swagger-flattened.json
@@ -19147,15 +19147,11 @@
                             "properties": {
                               "Code": {
                                 "description": "Specifies the family within a domain.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               },
                               "SubCode": {
                                 "description": "Specifies the sub-product family within a specific family.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               }
                             }
                           },
@@ -20273,15 +20269,11 @@
                             "properties": {
                               "Code": {
                                 "description": "Specifies the family within a domain.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               },
                               "SubCode": {
                                 "description": "Specifies the sub-product family within a specific family.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               }
                             }
                           },
@@ -35202,15 +35194,11 @@
                             "properties": {
                               "Code": {
                                 "description": "Specifies the family within a domain.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               },
                               "SubCode": {
                                 "description": "Specifies the sub-product family within a specific family.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               }
                             }
                           },
diff --git a/pkg/schema/spec/v3.1.4/account-info-swagger-flattened.json b/pkg/schema/spec/v3.1.4/account-info-swagger-flattened.json
index 0736b453..80d63b62 100644
--- a/pkg/schema/spec/v3.1.4/account-info-swagger-flattened.json
+++ b/pkg/schema/spec/v3.1.4/account-info-swagger-flattened.json
@@ -19257,15 +19257,11 @@
                               "properties": {
                                 "Code": {
                                   "description": "Specifies the family within a domain.",
-                                  "type": "string",
-                                  "minLength": 1,
-                                  "maxLength": 4
+                                  "type": "string"
                                 },
                                 "SubCode": {
                                   "description": "Specifies the sub-product family within a specific family.",
-                                  "type": "string",
-                                  "minLength": 1,
-                                  "maxLength": 4
+                                  "type": "string"
                                 }
                               }
                             },
@@ -20395,15 +20391,11 @@
                               "properties": {
                                 "Code": {
                                   "description": "Specifies the family within a domain.",
-                                  "type": "string",
-                                  "minLength": 1,
-                                  "maxLength": 4
+                                  "type": "string"
                                 },
                                 "SubCode": {
                                   "description": "Specifies the sub-product family within a specific family.",
-                                  "type": "string",
-                                  "minLength": 1,
-                                  "maxLength": 4
+                                  "type": "string"
                                 }
                               }
                             },
@@ -35407,15 +35399,11 @@
                               "properties": {
                                 "Code": {
                                   "description": "Specifies the family within a domain.",
-                                  "type": "string",
-                                  "minLength": 1,
-                                  "maxLength": 4
+                                  "type": "string"
                                 },
                                 "SubCode": {
                                   "description": "Specifies the sub-product family within a specific family.",
-                                  "type": "string",
-                                  "minLength": 1,
-                                  "maxLength": 4
+                                  "type": "string"
                                 }
                               }
                             },
diff --git a/pkg/schema/spec/v3.1.5/account-info-swagger-flattened.json b/pkg/schema/spec/v3.1.5/account-info-swagger-flattened.json
index dede8c5d..4c331eb8 100644
--- a/pkg/schema/spec/v3.1.5/account-info-swagger-flattened.json
+++ b/pkg/schema/spec/v3.1.5/account-info-swagger-flattened.json
@@ -19304,15 +19304,11 @@
                             "properties": {
                               "Code": {
                                 "description": "Specifies the family within a domain.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               },
                               "SubCode": {
                                 "description": "Specifies the sub-product family within a specific family.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               }
                             }
                           },
@@ -20451,15 +20447,11 @@
                             "properties": {
                               "Code": {
                                 "description": "Specifies the family within a domain.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               },
                               "SubCode": {
                                 "description": "Specifies the sub-product family within a specific family.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               }
                             }
                           },
@@ -35511,15 +35503,12 @@
                             "properties": {
                               "Code": {
                                 "description": "Specifies the family within a domain.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
                               },
                               "SubCode": {
                                 "description": "Specifies the sub-product family within a specific family.",
-                                "type": "string",
-                                "minLength": 1,
-                                "maxLength": 4
+                                "type": "string"
+
                               }
                             }
                           },
-- 
2.17.1
```