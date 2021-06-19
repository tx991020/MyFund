
select * from {{.table}} where c7 !='--' and c6 !='--' and c6::float > {{.day}}  and c7::float < {{.month}} order by c6::float desc limit 100