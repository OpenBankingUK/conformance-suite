defmodule Compliance.ValidationRuns.EndpointReportTest do
  @moduledoc """
  Tests for EndPointReport.
  """
  use ExUnit.Case, async: false
  alias Compliance.ValidationRuns.EndpointReport

  describe "EndpointReport" do
    test "new struct defaults" do
      report = %EndpointReport{path: 'http://example.com/api'}
      assert report.path == 'http://example.com/api'
      assert report.total_calls == 0
      assert report.failed_calls == 0
      assert report.failures == []
    end
  end
end
