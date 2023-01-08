# disco store queries 

> Note, `disc` tables are partitioned by day, for best performance use the `*_day` views which automatically scope the query to last day. For queries over time, use the pseudo-column named `_PARTITIONDATE` in your where clause to limit the query scope (https://cloud.google.com/bigquery/docs/querying-partitioned-tables)

## Most common licenses

```sql
SELECT
  name,
  count(*) num
FROM `cloudy-demos.disco.licenses_day`
GROUP BY 1
ORDER BY 2 desc
LIMIT 10
```

## Most common packages 

```sql
SELECT DISTINCT
  package, 
  count(*) num
FROM `cloudy-demos.disco.packages_day`
GROUP BY 1 
ORDER BY 2 DESC
```

## Unique CVEs for image 

```sql
SELECT DISTINCT
  cve,
  severity,
  package,
  version,
FROM `cloudy-demos.disco.vulnerabilities_day` 
WHERE image = 'us-west1-docker.pkg.dev/cloudy-demos/disco/disco'
ORDER BY 1 DESC
```

## Images with most vulnerabilities 

```sql
SELECT
  image, 
  SUM(lows) lows,
  SUM(mediums) mediums,
  SUM(highs) highs,
  SUM(criticals) criticals,
  SUM(unknowns) unknowns,
  COUNT(*) total
FROM (
  SELECT
    image, 
    CASE WHEN severity = 'LOW' THEN 1 ELSE 0 END lows,
    CASE WHEN severity = 'MEDIUM' THEN 1 ELSE 0 END mediums,
    CASE WHEN severity = 'HIGH' THEN 1 ELSE 0 END highs,
    CASE WHEN severity = 'CRITICAL' THEN 1 ELSE 0 END criticals,
    CASE WHEN severity = 'UNKNOWN' THEN 1 ELSE 0 END unknowns,
  FROM `cloudy-demos.disco.vulnerabilities_day`
) dt 
GROUP BY 1 
ORDER BY 7 DESC
LIMIT 10
```