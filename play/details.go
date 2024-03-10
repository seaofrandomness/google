package play

import (
   "154.pages.dev/encoding"
   "154.pages.dev/protobuf"
   "errors"
   "fmt"
   "io"
   "net/http"
)

func (d Details) File() (func() (uint64, bool), bool) {
   if v, ok := d.m.Get(13); ok {
      if v, ok := v.Get(1); ok {
         iterate := v.Iterate(17)
         return func() (uint64, bool) {
            if v, ok := iterate(); ok {
               if v, ok := v.GetVarint(1); ok {
                  return uint64(v), true
               }
            }
            return 0, false
         }, true
      }
   }
   return nil, false
}

func (d Details) String() string {
   var b []byte
   b = append(b, "downloads:"...)
   if v, ok := d.Downloads(); ok {
      b = fmt.Append(b, " ", encoding.Cardinal(v))
   }
   b = append(b, "\nfiles:"...)
   if iterate, ok := d.File(); ok {
      for {
         file, ok := iterate()
         if !ok {
            break
         }
         if file >= 1 {
            b = append(b, " OBB"...)
         } else {
            b = append(b, " APK"...)
         }
      }
   }
   b = append(b, "\nname:"...)
   if v, ok := d.Name(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\noffered by:"...)
   if v, ok := d.OfferedBy(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nprice:"...)
   if v, ok := d.Price(); ok {
      b = fmt.Append(b, " ", v)
   }
   if v, ok := d.PriceCurrency(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nrequires:"...)
   if v, ok := d.Requires(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nsize:"...)
   if v, ok := d.Size(); ok {
      b = fmt.Append(b, " ", encoding.Size(v))
   }
   b = append(b, "\nupdated on:"...)
   if v, ok := d.UpdatedOn(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nversion code:"...)
   if v, ok := d.VersionCode(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nversion name:"...)
   if v, ok := d.VersionName(); ok {
      b = fmt.Append(b, " ", v)
   }
   return string(b)
}

type Details struct {
   Checkin Checkin
   Token AccessToken
   m protobuf.Message
}

func (d Details) Name() (string, bool) {
   if v, ok := d.m.GetBytes(5); ok {
      return string(v), true
   }
   return "", false
}

func (d Details) OfferedBy() (string, bool) {
   if v, ok := d.m.GetBytes(6); ok {
      return string(v), true
   }
   return "", false
}

// developer.android.com/guide/topics/manifest/manifest-element
func (d Details) VersionCode() (uint64, bool) {
   d.m, _ = d.m.Get(13)
   d.m, _ = d.m.Get(1)
   if v, ok := d.m.GetVarint(3); ok {
      return uint64(v), true
   }
   return 0, false
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Downloads() (uint64, bool) {
   d.m, _ = d.m.Get(13)
   d.m, _ = d.m.Get(1)
   if v, ok := d.m.GetVarint(70); ok {
      return uint64(v), true
   }
   return 0, false
}

func (d Details) Price() (float64, bool) {
   d.m, _ = d.m.Get(8)
   if v, ok := d.m.GetVarint(1); ok {
      return float64(v) / 1_000_000, true
   }
   return 0, false
}

func (d Details) Size() (uint64, bool) {
   d.m, _ = d.m.Get(13)
   d.m, _ = d.m.Get(1)
   if v, ok := d.m.GetVarint(9); ok {
      return uint64(v), true
   }
   return 0, false
}

func (d Details) UpdatedOn() (string, bool) {
   d.m, _ = d.m.Get(13)
   d.m, _ = d.m.Get(1)
   if v, ok := d.m.GetBytes(16); ok {
      return string(v), true
   }
   return "", false
}

func (d Details) PriceCurrency() (string, bool) {
   d.m, _ = d.m.Get(8)
   if v, ok := d.m.GetBytes(2); ok {
      return string(v), true
   }
   return "", false
}

func (d Details) VersionName() (string, bool) {
   d.m, _ = d.m.Get(13)
   d.m, _ = d.m.Get(1)
   if v, ok := d.m.GetBytes(4); ok {
      return string(v), true
   }
   return "", false
}

func (d Details) Requires() (string, bool) {
   d.m, _ = d.m.Get(13)
   d.m, _ = d.m.Get(1)
   d.m, _ = d.m.Get(82)
   d.m, _ = d.m.Get(1)
   if v, ok := d.m.GetBytes(1); ok {
      return string(v), true
   }
   return "", false
}

func (d *Details) Details(app string, single bool) error {
   req, err := http.NewRequest("GET", "https://android.clients.google.com", nil)
   if err != nil {
      return err
   }
   req.URL.Path = "/fdfe/details"
   req.URL.RawQuery = "doc=" + app
   authorization(req, d.Token)
   user_agent(req, single)
   if err := x_dfe_device_id(req, d.Checkin); err != nil {
      return err
   }
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errors.New(res.Status)
   }
   data, err := io.ReadAll(res.Body)
   if err != nil {
      return err
   }
   if err := d.m.Consume(data); err != nil {
      return err
   }
   d.m, _ = d.m.Get(1)
   d.m, _ = d.m.Get(2)
   d.m, _ = d.m.Get(4)
   return nil
}
