package snap

import "testing"

func TestLatestRev(t *testing.T) {
	revisions := `Rev.    Uploaded              Arch     Version         Channels
1414    2020-01-21T22:46:23Z  ppc64el  1.15.9          1.15/stable*, 1.15/candidate*, 1.15/beta*, 1.15/edge*
1413    2020-01-21T22:36:55Z  ppc64el  1.17.2          1.17/candidate*, 1.17/beta*, 1.17/edge*
1412    2020-01-21T22:36:49Z  arm64    1.15.9          1.15/stable*, 1.15/candidate*, 1.15/beta*, 1.15/edge*
1411    2020-01-21T22:36:36Z  s390x    1.15.9          1.15/stable*, 1.15/candidate*, 1.15/beta*, 1.15/edge*
1410    2020-01-21T22:29:02Z  arm64    1.17.2          1.17/candidate*, 1.17/beta*, 1.17/edge*
1409    2020-01-21T22:23:16Z  amd64    1.15.9          1.15/stable*, 1.15/candidate*, 1.15/beta*, 1.15/edge*
1408    2020-01-21T22:19:18Z  amd64    1.17.2          1.17/candidate*, 1.17/beta*, 1.17/edge*, 1.17/stable*, stable*, edge*, beta*, candidate*
1407    2020-01-21T22:17:44Z  s390x    1.17.2          1.17/candidate*, 1.17/beta*, 1.17/edge*
1406    2020-01-17T14:07:32Z  ppc64el  1.15.8          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1405    2020-01-17T14:01:17Z  arm64    1.15.8          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1404    2020-01-17T13:58:48Z  arm64    1.15.8          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1403    2020-01-17T13:58:36Z  s390x    1.15.8          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1402    2020-01-17T13:57:27Z  amd64    1.15.8          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1401    2020-01-17T13:56:49Z  ppc64el  1.15.8          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1400    2020-01-17T13:52:54Z  s390x    1.15.8          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1399    2020-01-17T13:52:49Z  amd64    1.15.8          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1398    2020-01-15T14:41:13Z  arm64    1.17.1          1.17/candidate, 1.17/beta, 1.17/edge
1397    2020-01-15T14:29:23Z  arm64    1.17.1          1.17/candidate, 1.17/beta, 1.17/edge
1396    2020-01-15T14:27:20Z  ppc64el  1.17.1          1.17/candidate, 1.17/beta, 1.17/edge
1395    2020-01-15T14:24:31Z  amd64    1.17.1          -
1394    2020-01-15T14:24:20Z  amd64    1.17.1          1.17/candidate, 1.17/beta, 1.17/edge, 1.17/stable, edge, candidate, beta, stable
1393    2020-01-15T14:21:39Z  ppc64el  1.17.1          1.17/candidate, 1.17/beta, 1.17/edge
1392    2020-01-15T14:20:37Z  amd64    1.17.1          1.17/candidate, 1.17/beta, 1.17/edge
1391    2019-12-18T15:00:16Z  s390x    1.18.0-alpha.1  1.18/edge*
1390    2019-12-18T14:51:39Z  arm64    1.18.0-alpha.1  1.18/edge*
1389    2019-12-18T14:49:13Z  amd64    1.18.0-alpha.1  1.18/edge*
1388    2019-12-18T14:46:27Z  ppc64el  1.18.0-alpha.1  1.18/edge*
1387    2019-12-16T18:03:19Z  s390x    1.15.7          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1386    2019-12-16T17:58:50Z  s390x    1.16.4          1.16/stable*, 1.16/candidate*, 1.16/beta*, 1.16/edge*
1385    2019-12-16T16:48:33Z  s390x    1.14.10         1.14/stable*, 1.14/candidate*, 1.14/beta*, 1.14/edge*
1384    2019-12-16T16:47:45Z  arm64    1.15.7          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1383    2019-12-16T16:47:43Z  ppc64el  1.15.7          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1382    2019-12-16T16:46:47Z  ppc64el  1.16.4          1.16/stable*, 1.16/candidate*, 1.16/beta*, 1.16/edge*
1381    2019-12-16T16:45:45Z  arm64    1.16.4          1.16/stable*, 1.16/candidate*, 1.16/beta*, 1.16/edge*
1380    2019-12-16T16:38:27Z  amd64    1.16.4          1.16/stable*, 1.16/candidate*, 1.16/beta*, 1.16/edge*
1379    2019-12-16T16:37:28Z  arm64    1.14.10         1.14/stable*, 1.14/candidate*, 1.14/beta*, 1.14/edge*
1378    2019-12-16T16:35:42Z  amd64    1.15.7          1.15/stable, 1.15/candidate, 1.15/beta, 1.15/edge
1377    2019-12-16T16:32:10Z  amd64    1.14.10         1.14/stable*, 1.14/candidate*, 1.14/beta*, 1.14/edge*
1376    2019-12-16T16:31:18Z  ppc64el  1.14.10         1.14/stable*, 1.14/candidate*, 1.14/beta*, 1.14/edge*
1375    2019-12-10T03:25:31Z  arm64    1.17.0          1.17/stable*, 1.17/candidate, 1.17/beta, 1.17/edge
1374    2019-12-10T03:08:33Z  ppc64el  1.17.0          1.17/stable*, 1.17/candidate, 1.17/beta, 1.17/edge
1373    2019-12-10T03:08:19Z  amd64    1.17.0          1.17/stable, 1.17/candidate, 1.17/beta, 1.17/edge, edge, beta, candidate, stable
1372    2019-12-10T03:02:13Z  s390x    1.17.0          1.17/stable*, 1.17/candidate, 1.17/beta, 1.17/edge
1371    2019-12-04T17:24:14Z  arm64    1.17.0-rc.2     1.17/stable, 1.17/candidate, 1.17/beta, 1.17/edge
1370    2019-12-04T17:19:35Z  amd64    1.17.0-rc.2     -
1369    2019-12-04T17:17:14Z  ppc64el  1.17.0-rc.2     1.17/stable, 1.17/candidate, 1.17/beta, 1.17/edge
1368    2019-12-04T17:15:26Z  s390x    1.17.0-rc.2     -
1367    2019-12-04T17:13:20Z  s390x    1.17.0-rc.2     1.17/stable, 1.17/candidate, 1.17/beta, 1.17/edge
1366    2019-12-04T17:12:14Z  amd64    1.17.0-rc.2     1.17/stable, 1.17/candidate, 1.17/beta, 1.17/edge`
	LatestRev("kubectl", "amd64", "1.17/edge", revisions)
}
