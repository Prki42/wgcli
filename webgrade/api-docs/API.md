# API

Base url: `http://drwebgrade.matf.bg.ac.rs/backend`

## Login (POST)

Endpoint: `/server/services/authService.php`

Body (`application/x-www-form-urlencoded`):

```
username=...
password=...
action=login
```

Pri uspešnom login-u vraća 200 status code i json sa podacima o korisniku

Pri neuspešnom login-u vraća 200 status code i `"null"` u telu.

## Logout (GET)

Endpoint: `/server/services/authService.php?action=logout`

## Send submition (POST)

Endpoint: `/server/services/submissionService.php`

Body (`multipart/form-data`) _(primer)_:

```
-----------------------------36553843923648209551851707999
Content-Disposition: form-data; name="action"

postSubmission
-----------------------------36553843923648209551851707999
Content-Disposition: form-data; name="problemId"

1
-----------------------------36553843923648209551851707999
Content-Disposition: form-data; name="graderId"

1
-----------------------------36553843923648209551851707999
Content-Disposition: form-data; name="graderName"

C
-----------------------------36553843923648209551851707999
Content-Disposition: form-data; name="file"; filename="kvadrat_kub.c"
Content-Type: text/plain

#include<stdlib.h>
#include<stdio.h>

int main() {
    int n;
    printf("Unesite ceo broj: ");
    scanf("%d", &n);

    printf("Kvadrat: %d\nKub: %d\n", n*n, n*n*n);

    return EXIT_SUCCESS;
}
-----------------------------36553843923648209551851707999
Content-Disposition: form-data; name="template"

null
-----------------------------36553843923648209551851707999--
```

### Response

Primer:

```json
{
  "graderTaskId": "fec599d2-c08e-3d65-8458-d1f74f71fkb0",
  "submissionId": "223461"
}
```

## Check submission tests (GET)

Endpoint _(primer)_: `/server/services/submissionService.php?action=getTaskStatus&graderTaskId=fec599d2-c08e-3d65-8458-d1f74f71fkb0&submissionId=233157`

GET zahtevi se moraju slati na ovaj endpoint sve dok se ne dobije `"state": "finished"` u odgovoru i tek onda se (pretpostavljam) na drugom servisu evidentira procenat tačnih primera.

To jest, servis za proveru koda nikad neće sam da obavesti drugi servis o postignutom rezultatu sve dok se ne pošalje zahtev za proveru. Potreba za slanjem submissionId-a implicira da servis za proveru koda nema podatak o tome "za koj submission" on radi tu proveru i ako je u momentu slanja zahteva proces završen score će biti update-ovan onom submission-u čiji je id je poslat u zahtevu. Ovo je i eksperimentalno potvrđeno.

### Response

Primer 1:

```json
{
  "state": "running"
}
```

Primer 2:

```json
{
  "state": "running",
  "tests": [
    {
      "status": "finished",
      "code": 0,
      "output": 1,
      "time": 0,
      "wall_time": 0,
      "memory": 262144
    }
  ]
}
```

Primer 3:

```json
{
  "state": "finished",
  "tests": [
    {
      "status": "finished",
      "code": 0,
      "output": 1,
      "time": 0,
      "wall_time": 0,
      "memory": 262144
    },
    {
      "status": "finished",
      "code": 0,
      "output": 1,
      "time": 0,
      "wall_time": 0,
      "memory": 262144
    },
    {
      "status": "finished",
      "code": 0,
      "output": 1,
      "time": 0,
      "wall_time": 0,
      "memory": 262144
    }
  ]
}
```

## Get submition details (GET)

Endpoint _(primer)_: `/server/services/submissionService.php?action=getSubmissionDetails&submissionId=233157`

### Response

Primer:

```json
{
  "score": "100",
  "timestamp": "2023-01-07 22:04:59",
  "requestedReview": "0",
  "requestText": null,
  "sourceFile": "#include<stdlib.h>\n#include<stdio.h>\n\nint main() {\n    int n;\n    printf(\"Unesite ceo broj: \");\n    scanf(\"%d\", &n);\n\n    printf(\"Kvadrat: %d\\nKub: %d\\n\", n*n, n*n*n);\n\n    return EXIT_SUCCESS;\n}",
  "codeReview": null,
  "reviewer": null,
  "submissionId": "233157",
  "reviewerAvatar": null,
  "reviewerId": null,
  "seenReview": "0",
  "extension": "2048/1673125499kvadrat_kub.c"
}
```

## Get problem details (GET)

Endpoint _(primer)_: `/server/services/studentService.php?action=getProblemDetailsComponentData&courseId=1&topicId=1&problemId=1`

`courseId` i `topicId` mogu biti postavljeni na `"null"` što dovodi do toga da će u odgvoru `breadcrumb.course` i `breadcrumb.topic` takođe biti `null`. Slično, `courseId` i `topicId` se mogu proizvoljno staviti bilo koje id-eve (iako, npr, problem nije u tom kursu na primer) i biće vraćeni podaci o tom kursu i/ili topic-u.

### Response

Primer:

````json
{
  "breadcrumbData": {
    "course": { "name": "Programiranje 1 (Matematika)" },
    "topic": "Tema 4: Uvod u programski jezik C",
    "problem": { "name": "Kvadrat i kub" }
  },
  "authData": {
    "username": "...",
    "userId": 2048,
    "email": "...",
    "roleName": "student",
    "isVerified": 1
  },
  "problemSubmissions": [
    {
      "score": 0,
      "timestamp": "2023-01-01 16:34:26",
      "requestedReview": 0,
      "requestText": null,
      "sourceFile": "#include<stdlib.h>\n#include<stdio.h>\n\nint main() {\n    int n;\n    scanf(\"%d\", &n);\n\n    printf(\"%d\\n%d\\n\", n*n, n*n*n);\n\n    return EXIT_SUCCESS;\n}",
      "codeReview": null,
      "reviewer": null,
      "submissionId": 231868,
      "reviewerAvatar": null,
      "reviewerId": null,
      "seenReview": null,
      "extension": "2048/1672587266kvadrat_kub.c"
    },
    {
      "score": 100,
      "timestamp": "2023-01-01 16:35:14",
      "requestedReview": 0,
      "requestText": null,
      "sourceFile": "#include<stdlib.h>\n#include<stdio.h>\n\nint main() {\n    int n;\n    printf(\"Unesite ceo broj: \");\n    scanf(\"%d\", &n);\n\n    printf(\"Kvadrat: %d\\nKub: %d\\n\", n*n, n*n*n);\n\n    return EXIT_SUCCESS;\n}",
      "codeReview": null,
      "reviewer": null,
      "submissionId": 231869,
      "reviewerAvatar": null,
      "reviewerId": null,
      "seenReview": null,
      "extension": "2048/1672587314kvadrat_kub.c"
    }
  ],
  "problemDetails": {
    "problemDetails": {
      "name": "Kvadrat i kub",
      "text": "Napisati program koji omogu\u0107ava korisniku da unese ceo broj, a zatim ispisuje njegov kvadrat i kub.\n\nInterakcija sa programom:\n```\nUnesite ceo broj: 4\nKvadrat:16\nKub: 64\n```",
      "timeLimit": 1000,
      "memoryLimit": 67108864,
      "testCasesNum": 3,
      "templateFileContent": null
    },
    "graders": [{ "id": 1, "name": "C" }]
  },
  "professors": [
    {
      "idUser": 279,
      "avatar": "...",
      "firstName": "...",
      "lastName": "..."
    },
    {
      "idUser": 863,
      "avatar": "wg_admin",
      "firstName": "WG",
      "lastName": "Admin"
    }
  ],
  "reports": [],
  "canSubmitReport": true,
  "companyData": null
}
````
