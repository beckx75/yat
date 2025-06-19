var _____WB$wombat$assign$function_____ = function(name) {return (self._wb_wombat && self._wb_wombat.local_init && self._wb_wombat.local_init(name)) || self[name]; };
if (!self.__WB_pmw) { self.__WB_pmw = function(obj) { this.__WB_source = obj; return this; } }
{
  let window = _____WB$wombat$assign$function_____("window");
  let self = _____WB$wombat$assign$function_____("self");
  let document = _____WB$wombat$assign$function_____("document");
  let location = _____WB$wombat$assign$function_____("location");
  let top = _____WB$wombat$assign$function_____("top");
  let parent = _____WB$wombat$assign$function_____("parent");
  let frames = _____WB$wombat$assign$function_____("frames");
  let opener = _____WB$wombat$assign$function_____("opener");

function toggleMenu (objID)
{
    if (document.getElementById (objID).style.display != "block")
   {
        document.getElementById (objID).style.display = "block";
        document.getElementById ('togglelink').innerHTML = "[ fewer options ]";
    }
    else
    {
        document.getElementById (objID).style.display = "none";
        document.getElementById ('togglelink').innerHTML = "[ more options ]";
     }
}


}
/*
     FILE ARCHIVED ON 05:58:11 Mar 28, 2020 AND RETRIEVED FROM THE
     INTERNET ARCHIVE ON 05:49:41 Nov 03, 2023.
     JAVASCRIPT APPENDED BY WAYBACK MACHINE, COPYRIGHT INTERNET ARCHIVE.

     ALL OTHER CONTENT MAY ALSO BE PROTECTED BY COPYRIGHT (17 U.S.C.
     SECTION 108(a)(3)).
*/
/*
playback timings (ms):
  captures_list: 67.462
  exclusion.robots: 0.072
  exclusion.robots.policy: 0.063
  cdx.remote: 0.06
  esindex: 0.008
  LoadShardBlock: 40.178 (3)
  PetaboxLoader3.datanode: 90.312 (5)
  load_resource: 232.765 (2)
  PetaboxLoader3.resolve: 108.82 (2)
*/