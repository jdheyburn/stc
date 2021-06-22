package repository

// TODO tidy this up or use gorm (RIP)
var nlcs_query = `with 
this_loc as (
	select nlc, fare_group, zone_no from location where crs = ?
),
this_group_members as (
	select group_uic_code, member_crs_code from location_group_member lgm where member_crs_code = ?
),
this_group_members_loc as (
	select location.nlc from this_group_members left join location on this_group_members.group_uic_code = location.uic
),
nlcs as (
	-- query 1 - this location NLC
	Select nlc from this_loc
	union
	-- query 2 - clustered location NLC
	select cluster_id as nlc from this_loc left join station_cluster on this_loc.nlc = station_cluster.cluster_nlc where cluster_id is not null
	union
	-- query 3 - NLCs that exist in other groups
	select nlc from this_group_members_loc
	union
	-- query 4 - clustered locations from group members
	select cluster_id as nlc from this_group_members_loc left join station_cluster on this_group_members_loc.nlc = station_cluster.cluster_nlc where station_cluster.cluster_nlc is not null
	union
	-- query 5 - this location fare group
	select fare_group as nlc from this_loc
	union
	-- query 6 - clustered locations from fare group
	select cluster_id as nlc from this_loc left join station_cluster on this_loc.fare_group = station_cluster.cluster_nlc  where cluster_id is not null
	union
	-- query 7 - (mine) lookup against zone group
	select zone_no as nlc from this_loc where zone_no is not null
)
select distinct(nlc) from nlcs`

// Apologies in advance
var fares_query = `select distinct
flow.origin_code
, (select description from location where nlc = flow.origin_code and location.start_date <= CURDATE() and location.end_date > CURDATE()) as origin_name
, flow.destination_code
, (select description from location where nlc = flow.destination_code and location.start_date <= CURDATE() and location.end_date > CURDATE()) as destination_name
, flow.route_code
, route.description as route_desc
, route.aaa_desc as route_aaa_desc
, flow.status_code
, flow.usage_code
, flow.toc
, fare.flow_id
, fare.id as fare_id
, fare.ticket_code
,ticket_type.description as ticket_desc
,ticket_type.tkt_class as ticket_class
,ticket_type.tkt_type as ticket_type
, fare.fare as adult_fare
, cast((fare.fare / 2) as UNSIGNED) as child_fare
, fare.restriction_code
, restriction_header.description as restriction_desc
from flow 
left join route on flow.route_code = route.route_code 
left join fare on flow.flow_id = fare.flow_id 
left join ticket_type on fare.ticket_code = ticket_type.ticket_code 
left join restriction_header on fare.restriction_code = restriction_header.restriction_code
where 
flow.start_date <= CURDATE() and flow.end_date > CURDATE() 
AND route.start_date <= CURDATE() and route.end_date > CURDATE()
AND ticket_type.start_date <= CURDATE() and ticket_type.end_date > CURDATE()
AND ticket_type.tkt_class = '2'
AND flow.route_code IN ('00000', '01000') -- default to any permitted routes for now
AND (
	(origin_code IN ? and destination_code in ?) 
	OR
	(origin_code IN ? and destination_code in ? and direction = 'R')
)
order by fare asc
`

var nfo_query = `select distinct
ndfo.origin_code
, (select description from location where nlc = ndfo.origin_code and location.start_date <= CURDATE() and location.end_date > CURDATE()) as origin_name
, ndfo.destination_code
, (select description from location where nlc = ndfo.destination_code and location.start_date <= CURDATE() and location.end_date > CURDATE()) as destination_name
, ndfo.route_code 
, route.description as route_desc
, route.aaa_desc as route_aaa_desc
, NULL as status_code
, NULL as usage_code
, NULL as toc
, NULL as flow_id
, NULL as fare_id
, ndfo.ticket_code 
,ticket_type.description as ticket_desc
,ticket_type.tkt_class as ticket_class
,ticket_type.tkt_type as ticket_type
, ndfo.adult_fare
, ndfo.child_fare
, ndfo.restriction_code 
, restriction_header.description as restriction_desc
from non_derivable_fare_override ndfo
left join route on ndfo.route_code = route.route_code 
left join ticket_type on ndfo.ticket_code = ticket_type.ticket_code 
left join restriction_header on ndfo.restriction_code = restriction_header.restriction_code
where 
ndfo.start_date <= CURDATE() and ndfo.end_date > CURDATE() 
and railcard_code = ''
and origin_code IN ?
and destination_code in ?
`
