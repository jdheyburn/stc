-- Step 1: Get NLC from CRS codes

select uic, start_date, end_date, nlc, description, crs, fare_group from location where crs = 'SNR' and start_date <= CURDATE() and end_date > CURDATE();

|uic    |start_date|end_date  |nlc |description     |crs|fare_group|
|-------|----------|----------|----|----------------|---|----------|
|7054330|2020-09-09|2999-12-31|5433|SANDERSTEAD     |SNR|5433      |

select uic, start_date, end_date, nlc, description, crs, fare_group from location where crs = 'EGR' and start_date <= CURDATE() and end_date > CURDATE();

|uic    |start_date|end_date  |nlc |description     |crs|fare_group|
|-------|----------|----------|----|----------------|---|----------|
|7054860|2020-09-09|2999-12-31|5486|EAST GRINSTEAD  |EGR|5486      |


-- Step 2: Get flows between the two NLCs

-- EGR -> SNR

SELECT flow.flow_id,flow.origin_code,flow.destination_code,flow.direction,flow.start_date,flow.end_date,flow.route_code,route.description as route_desc 
FROM `flow`
LEFT JOIN route on flow.route_code = route.route_code 
WHERE flow.origin_code = 5486 AND flow.destination_code = 5433 AND flow.start_date <= CURDATE() AND flow.end_date > CURDATE() AND route.start_date <= CURDATE() AND route.end_date > CURDATE()


|flow_id|origin_code|destination_code|direction|start_date|end_date  |route_code|route_desc      |
|-------|-----------|----------------|---------|----------|----------|----------|----------------|
|137711 |5486       |5433            |S        |2020-01-02|2999-12-31|01000     |.               |


-- SNR -> EGR

SELECT flow.flow_id,flow.origin_code,flow.destination_code,flow.direction,flow.start_date,flow.end_date,flow.route_code,route.description as route_desc 
FROM `flow`
LEFT JOIN route on flow.route_code = route.route_code 
WHERE flow.origin_code = 5433 AND flow.destination_code = 5486 AND flow.start_date <= CURDATE() AND flow.end_date > CURDATE() AND route.start_date <= CURDATE() AND route.end_date > CURDATE()



|flow_id|origin_code|destination_code|direction|start_date|end_date  |route_code|route_desc      |
|-------|-----------|----------------|---------|----------|----------|----------|----------------|
|136991 |5433       |5486            |S        |2020-01-02|2999-12-31|01000     |.               |


-- Step 3: Get fares for the given flow ID

-- I've listed both out here to see if there is any difference between the flow fares (doesn't look like it)

-- EGR -> SNR
SELECT fare.id,fare.flow_id,fare.fare,fare.ticket_code,ticket_type.description as ticket_description,ticket_type.tkt_class as ticket_class,ticket_type.tkt_type as ticket_type,
fare.restriction_code, restriction_header.description as restriction_desc, restriction_header.desc_out as restriction_desc_out, restriction_header.desc_ret as restriction_desc_rtn
FROM `fare` 
LEFT JOIN ticket_type on fare.ticket_code = ticket_type.ticket_code 
LEFT JOIN restriction_header on fare.restriction_code = restriction_header.restriction_code 
WHERE fare.flow_id IN (137711) AND ticket_type.start_date <= CURDATE() AND ticket_type.end_date > CURDATE()


|id         |flow_id|fare    |ticket_code|ticket_description|ticket_class|ticket_type|restriction_code|restriction_desc              |restriction_desc_out                              |restriction_desc_rtn                              |
|-----------|-------|--------|-----------|------------------|------------|-----------|----------------|------------------------------|--------------------------------------------------|--------------------------------------------------|
|1064475    |137711 |1430    |0AJ        |SMART FCR         |1           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1064475    |137711 |1430    |0AJ        |SMART FCR         |1           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1064476    |137711 |950     |0AK        |SMART CDR         |2           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1064476    |137711 |950     |0AK        |SMART CDR         |2           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1064483    |137711 |950     |CDR        |OFF-PEAK DAY R    |2           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1064483    |137711 |950     |CDR        |OFF-PEAK DAY R    |2           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1064484    |137711 |1430    |FCR        |OFF-PEAK DAY 1R   |1           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1064484    |137711 |1430    |FCR        |OFF-PEAK DAY 1R   |1           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1064471    |137711 |2430    |0AC        |SMART FDR         |1           |R          |                |                              |                                                  |                                                  |
|1064472    |137711 |1280    |0AD        |SMART FDS         |1           |S          |                |                              |                                                  |                                                  |
|1064473    |137711 |1620    |0AE        |SMART SDR         |2           |R          |                |                              |                                                  |                                                  |
|1064474    |137711 |850     |0AF        |SMART SDS         |2           |S          |                |                              |                                                  |                                                  |
|1064477    |137711 |5300    |0AQ        |SMART 7DS         |2           |N          |                |                              |                                                  |                                                  |
|1064478    |137711 |7950    |0AR        |SMART 7DF         |1           |N          |                |                              |                                                  |                                                  |
|1064479    |137711 |5300    |0AS        |SMART PSS         |2           |N          |                |                              |                                                  |                                                  |
|1064480    |137711 |7950    |0AT        |SMART PSF         |1           |N          |                |                              |                                                  |                                                  |
|1064481    |137711 |7950    |7DF        |SEVEN DAY   1ST   |1           |N          |                |                              |                                                  |                                                  |
|1064482    |137711 |5300    |7DS        |SEVEN DAY   STD   |2           |N          |                |                              |                                                  |                                                  |
|1064485    |137711 |2430    |FDR        |ANYTIME DAY 1R    |1           |R          |                |                              |                                                  |                                                  |
|1064486    |137711 |1280    |FDS        |ANYTIME DAY 1S    |1           |S          |                |                              |                                                  |                                                  |
|1064487    |137711 |1620    |SDR        |ANYTIME DAY R     |2           |R          |                |                              |                                                  |                                                  |
|1064488    |137711 |850     |SDS        |ANYTIME DAY S     |2           |S          |                |                              |                                                  |                                                  |
|1064489    |137711 |400     |TKR        |CHILD FLTFARE R   |2           |R          |                |                              |                                                  |                                                  |
|1064490    |137711 |400     |TKS        |CHILD FLTFARE S   |2           |S          |                |                              |                                                  |                                                  |

-- SNR -> EGR
SELECT fare.id,fare.flow_id,fare.fare,fare.ticket_code,ticket_type.description as ticket_description,ticket_type.tkt_class as ticket_class,ticket_type.tkt_type as ticket_type,
fare.restriction_code, restriction_header.description as restriction_desc, restriction_header.desc_out as restriction_desc_out, restriction_header.desc_ret as restriction_desc_rtn
FROM `fare` 
LEFT JOIN ticket_type on fare.ticket_code = ticket_type.ticket_code 
LEFT JOIN restriction_header on fare.restriction_code = restriction_header.restriction_code 
WHERE fare.flow_id IN (136991) AND ticket_type.start_date <= CURDATE() AND ticket_type.end_date > CURDATE()

|id         |flow_id|fare    |ticket_code|ticket_description|ticket_class|ticket_type|restriction_code|restriction_desc              |restriction_desc_out                              |restriction_desc_rtn                              |
|-----------|-------|--------|-----------|------------------|------------|-----------|----------------|------------------------------|--------------------------------------------------|--------------------------------------------------|
|1058087    |136991 |1430    |0AJ        |SMART FCR         |1           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1058087    |136991 |1430    |0AJ        |SMART FCR         |1           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1058088    |136991 |950     |0AK        |SMART CDR         |2           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1058088    |136991 |950     |0AK        |SMART CDR         |2           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1058095    |136991 |950     |CDR        |OFF-PEAK DAY R    |2           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1058095    |136991 |950     |CDR        |OFF-PEAK DAY R    |2           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1058096    |136991 |1430    |FCR        |OFF-PEAK DAY 1R   |1           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1058096    |136991 |1430    |FCR        |OFF-PEAK DAY 1R   |1           |R          |B1              |OFF-PEAK                      |VALID AFTER 0929 MON-FRI                          |VALID AFTER 0929 MON-FRI                          |
|1058083    |136991 |1770    |0AC        |SMART FDR         |1           |R          |                |                              |                                                  |                                                  |
|1058084    |136991 |1320    |0AD        |SMART FDS         |1           |S          |                |                              |                                                  |                                                  |
|1058085    |136991 |1180    |0AE        |SMART SDR         |2           |R          |                |                              |                                                  |                                                  |
|1058086    |136991 |880     |0AF        |SMART SDS         |2           |S          |                |                              |                                                  |                                                  |
|1058089    |136991 |5300    |0AQ        |SMART 7DS         |2           |N          |                |                              |                                                  |                                                  |
|1058090    |136991 |7950    |0AR        |SMART 7DF         |1           |N          |                |                              |                                                  |                                                  |
|1058091    |136991 |5300    |0AS        |SMART PSS         |2           |N          |                |                              |                                                  |                                                  |
|1058092    |136991 |7950    |0AT        |SMART PSF         |1           |N          |                |                              |                                                  |                                                  |
|1058093    |136991 |7950    |7DF        |SEVEN DAY   1ST   |1           |N          |                |                              |                                                  |                                                  |
|1058094    |136991 |5300    |7DS        |SEVEN DAY   STD   |2           |N          |                |                              |                                                  |                                                  |
|1058097    |136991 |1770    |FDR        |ANYTIME DAY 1R    |1           |R          |                |                              |                                                  |                                                  |
|1058098    |136991 |1320    |FDS        |ANYTIME DAY 1S    |1           |S          |                |                              |                                                  |                                                  |
|1058099    |136991 |1180    |SDR        |ANYTIME DAY R     |2           |R          |                |                              |                                                  |                                                  |
|1058100    |136991 |880     |SDS        |ANYTIME DAY S     |2           |S          |                |                              |                                                  |                                                  |
|1058101    |136991 |400     |TKR        |CHILD FLTFARE R   |2           |R          |                |                              |                                                  |                                                  |
|1058102    |136991 |400     |TKS        |CHILD FLTFARE S   |2           |S          |                |                              |                                                  |                                                  |




-- Season tickets are calculated from 7DS for 2nd class and 7DF for 1st class
-- Steps 4 and 5 hidden; they pick out 7DS above and then perform the calculations





SELECT fare.id,fare.flow_id,fare.fare,fare.ticket_code,ticket_type.description as ticket_description,ticket_type.tkt_class as ticket_class,ticket_type.tkt_type as ticket_type,
fare.restriction_code, restriction_header.description as restriction_desc, restriction_header.desc_out as restriction_desc_out, restriction_header.desc_ret as restriction_desc_rtn
FROM `fare` 
LEFT JOIN ticket_type on fare.ticket_code = ticket_type.ticket_code 
LEFT JOIN restriction_header on fare.restriction_code = restriction_header.restriction_code 
WHERE fare.flow_id IN (137711) AND ticket_type.start_date <= CURDATE() AND ticket_type.end_date > CURDATE()


select * from restriction_header re where restriction_code = 'PG';
















select * from flow where origin_code = '1072' or destination_code = '1072';


with f as (
	select destination_code as c from flow where origin_code = '5148'
	union
	select origin_code as c from flow where destination_code = '5148'
)
select distinct f.c, location.description from f
inner join location on location.nlc = f.c
order by description ;

select uic, start_date, end_date, nlc, description, crs, fare_group from location where crs = 'SNR' and start_date <= CURDATE() and end_date > CURDATE();
-- SNR uic 7054330

select uic, start_date, end_date, nlc, description, crs, fare_group from location where uic = '7000320' and start_date <= CURDATE() and end_date > CURDATE();


select * from location_group where group_uic_code = '7010720';

select * from location_group_member where member_crs_code = 'VIC';

-- select location_group_member.group_uic_code, location_group.description as group_description
-- TODO alter below query so that it can be used for grouped and non-grouped stations

select 
location_group_member.member_uic_code
, location.nlc
, location.description
, location_group_member.member_crs_code
, location_group.description as group_description
, location_group_member.group_uic_code 
,location.fare_group
, location.start_date
, location.end_date
from location_group_member 
left join location_group on location_group_member.group_uic_code  = location_group.group_uic_code 
left join location on location_group.group_uic_code = location.uic 
where location_group_member.member_crs_code = 'VIC' 
and location_group_member.end_date > CURDATE()
and location_group.start_date <= CURDATE() and location_group.end_date > CURDATE()
and location.start_date <= CURDATE() and location.end_date > CURDATE();





select location.uic, location.start_date, location.end_date, location.nlc, location.description, location.crs, location.fare_group, location_group_member.group_uic_code, location_group_member.member_uic_code
from location 
left join location_group_member on location_group_member.member_crs_code = location.crs
where location.crs = 'LBG' and location.start_date <= CURDATE() and location.end_date > CURDATE() and location_group_member.end_date > CURDATE() and location_group_member.member_crs_code = 'LBG';



-- LBG nlc = 5148
-- LBG fare_group = 1072