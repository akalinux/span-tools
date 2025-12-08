package st
import "testing"
import "iter"


func TestAccumulateIterPull2(t *testing.T) {

  var next, stop = iter.Pull2(testDriver.NewSpanOverlapAccumulator().SliceIterFactory(&MultiSet))
  defer stop()
  var idx, res, hasnext = next()
  var count = -1
  for hasnext {
    count++
    if idx != count {
      t.Errorf("Failed, expected %d, got %d", count, idx)
      return
    }
    switch idx {
    case 0:
      {
        if res.GetBegin() != -1 || res.GetEnd() != 0 {
          t.Errorf("Invalid Range on set 0, expected: -1,0, got %d,%d", res.GetBegin(), res.GetEnd())
          return
        }
        if res.Contains != nil {
          t.Errorf("Exepcted Empty contains")
          return
        }
      }
    case 1:
      {
        if res.GetBegin() != 2 || res.GetEnd() != 2 {
          t.Errorf("Invalid Range on set 1, expected: 2,2, got %d,%d", res.GetBegin(), res.GetEnd())
          return
        }
        if res.GetContains() == nil {
          t.Errorf("Exepcted Non-Empty contains")
          return
        }
      }
    case 2:
      {
        if res.GetBegin() != 5 || res.GetEnd() != 6 {
          t.Errorf("Invalid Range on set 0, expected: 5,6... got %d,%d", res.GetBegin(), res.GetEnd())
          return
        }
        if res.Contains != nil {
          t.Errorf("Exepcted Empty contains")
          return
        }
      }
    case 3:
      {
        if res.GetBegin() != 9 || res.GetEnd() != 11 {
          t.Errorf("Invalid Range on set 0, expected: 9,11... got %d,%d", res.GetBegin(), res.GetEnd())
          return
        }
        if res.Contains != nil {
          t.Errorf("Exepcted Empty contains")
          return
        }
      }
    default:
      {
        t.Errorf("Got a range beyond 3, expected set to end at the, end is at: %d", idx)
        return
      }
    }
    idx, res, hasnext = next()
  }
  if count != 3 {
    t.Errorf("Expected 3 rows, got %d", count)
  }
}
