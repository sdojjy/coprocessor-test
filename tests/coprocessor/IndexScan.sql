use coprocessor;
SET time_zone='+00:00';

select * from t1 where str_unq_idx = 'ecczrvzje' or str_unq_idx = 'sttyhtbzjdeqbdqhdjrhzaxkrglyprylxvkuurffbgwbleiapcddfmxbgipbjvpfnzzwuugnsuyeuegzksabaslcnl';
select str_unq_idx from t1 where str_unq_idx is not null order by str_unq_idx;
select str_unq_idx from t1 where str_unq_idx > 'e' and str_unq_idx < 'f' order by id;
select * from t1 where str_unq_idx = 'tzephezkototsshyvqudwiphtbxlyyuialdsjmzbuuftateahsmhdqvbvbdpkpvxdewogzvtydhgnagtiibewlcligxkniffjhkmoujirfrlyvbx' or str_unq_idx < 'e' order by str_unq_idx;
select * from t1 where int_idx < 9 order by str_unq_idx;
select * from t1 where str_idx1_1 like 'd%' order by id;
select * from t1 where str_idx1_1 = 'nwfqlhhkkdijltmtziitvgcbfnbdnsdx' and str_idx1_2 like 'j%' order by id;
