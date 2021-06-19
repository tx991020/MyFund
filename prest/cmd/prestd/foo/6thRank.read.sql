
SELECT * from {{.table}} where c13 <> '--' and c9 != '--' and c13::float < 14  order by c9::float desc limit 60